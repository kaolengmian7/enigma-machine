package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Rotor 表示 Enigma 机器中的转子
type Rotor struct {
	wiring   string // 转子内部接线
	position int    // 当前位置
	notch    int    // 触发下一个转子转动的凹槽位置
}

// Reflector 表示反射器
type Reflector struct {
	wiring string
}

// Plugboard 表示接线板
type Plugboard struct {
	connections map[rune]rune
}

// EnigmaMachine 表示完整的 Enigma 加密机
type EnigmaMachine struct {
	rotors    [3]Rotor
	reflector Reflector
	plugboard Plugboard
}

// NewEnigmaMachine 创建一个新的 Enigma 加密机实例
func NewEnigmaMachine() *EnigmaMachine {
	// 创建三个转子，使用历史上真实的 Enigma I 转子接线方式
	rotors := [3]Rotor{
		{wiring: "EKMFLGDQVZNTOWYHXUSPAIBRCJ", position: 0, notch: 16}, // Rotor I
		{wiring: "AJDKSIRUXBLHWTMCQGZNPYFVOE", position: 0, notch: 4},  // Rotor II
		{wiring: "BDFHJLCPRTXVZNYEIWGAKMUSQO", position: 0, notch: 21}, // Rotor III
	}

	// 使用历史上的 B 反射器接线方式
	reflector := Reflector{wiring: "YRUHQSLDPXNGOKMIEBFZCWVJAT"}

	// 创建一个空的接线板
	plugboard := Plugboard{connections: make(map[rune]rune)}

	return &EnigmaMachine{
		rotors:    rotors,
		reflector: reflector,
		plugboard: plugboard,
	}
}

// SetPlugboard 设置接线板连接
func (e *EnigmaMachine) SetPlugboard(pairs []string) error {
	if len(pairs) > 10 {
		return fmt.Errorf("接线板最多支持10对连接")
	}

	e.plugboard.connections = make(map[rune]rune)
	for _, pair := range pairs {
		if len(pair) != 2 {
			return fmt.Errorf("无效的接线对: %s", pair)
		}
		a, b := rune(pair[0]), rune(pair[1])
		e.plugboard.connections[a] = b
		e.plugboard.connections[b] = a
	}
	return nil
}

// SetRotorPositions 设置转子初始位置
func (e *EnigmaMachine) SetRotorPositions(positions [3]int) {
	for i := 0; i < 3; i++ {
		e.rotors[i].position = positions[i] % 26
	}
}

// rotateRotors 转动转子
func (e *EnigmaMachine) rotateRotors() {
	// 第一个转子每次都转动
	e.rotors[0].position = (e.rotors[0].position + 1) % 26

	// 检查是否需要转动第二个转子
	if e.rotors[0].position == e.rotors[0].notch {
		e.rotors[1].position = (e.rotors[1].position + 1) % 26

		// 检查是否需要转动第三个转子
		if e.rotors[1].position == e.rotors[1].notch {
			e.rotors[2].position = (e.rotors[2].position + 1) % 26
		}
	}
}

// encryptLetter 加密单个字母
func (e *EnigmaMachine) encryptLetter(letter rune) rune {
	// 转动转子
	e.rotateRotors()

	// 通过接线板
	if mapped, ok := e.plugboard.connections[letter]; ok {
		letter = mapped
	}

	// 将字母转换为 0-25 的数字
	pos := int(letter - 'A')

	// 正向通过三个转子
	for i := 0; i < 3; i++ {
		pos = (pos + e.rotors[i].position) % 26
		pos = int(e.rotors[i].wiring[pos] - 'A')
		pos = (pos - e.rotors[i].position + 26) % 26
	}

	// 通过反射器
	pos = int(e.reflector.wiring[pos] - 'A')

	// 反向通过三个转子
	for i := 2; i >= 0; i-- {
		pos = (pos + e.rotors[i].position) % 26
		reversed := strings.IndexByte(e.rotors[i].wiring, byte(pos+'A'))
		pos = (reversed - e.rotors[i].position + 26) % 26
	}

	// 再次通过接线板
	result := rune(pos + 'A')
	if mapped, ok := e.plugboard.connections[result]; ok {
		result = mapped
	}

	return result
}

// Encrypt 加密消息
func (e *EnigmaMachine) Encrypt(message string) string {
	message = strings.ToUpper(message)
	var result strings.Builder

	for _, letter := range message {
		if letter >= 'A' && letter <= 'Z' {
			result.WriteRune(e.encryptLetter(letter))
		} else {
			result.WriteRune(letter)
		}
	}

	return result.String()
}

// 添加请求和响应的结构体
type EncryptRequest struct {
	Message   string   `json:"message"`
	Plugboard []string `json:"plugboard"`
	Positions [3]int   `json:"positions"`
}

type EncryptResponse struct {
	Result string `json:"result"`
}

// 处理加密请求的处理器
func handleEncrypt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	var req EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	// 创建新的 Enigma 实例
	enigma := NewEnigmaMachine()

	// 设置接线板
	if err := enigma.SetPlugboard(req.Plugboard); err != nil {
		http.Error(w, "接线板配置错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 设置转子位置
	enigma.SetRotorPositions(req.Positions)

	// 加密消息
	result := enigma.Encrypt(req.Message)

	// 返回结果
	response := EncryptResponse{Result: result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 修改 main 函数
func main() {
	// 设置路由
	http.HandleFunc("/api/encrypt", handleEncrypt)
	http.HandleFunc("/api/decrypt", handleEncrypt) // 解密使用相同的加密逻辑

	// 启动服务器
	port := ":8080"
	log.Printf("启动服务器在端口 %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
