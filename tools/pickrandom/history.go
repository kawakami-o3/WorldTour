package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// History は選択済み国の履歴
type History struct {
	Picked []string `json:"picked"`
}

// historyPath はプロジェクトルートからの相対パスで履歴ファイルを返す
func historyPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	// bin/pickrandom から2つ上がプロジェクトルート
	root := filepath.Dir(filepath.Dir(exe))
	return filepath.Join(root, "data", "history.json"), nil
}

// LoadHistory は履歴ファイルを読み込む。ファイルがなければ空の履歴を返す。
func LoadHistory() (*History, error) {
	path, err := historyPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{Picked: []string{}}, nil
		}
		return nil, err
	}

	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// Save は履歴をファイルに保存する
func (h *History) Save() error {
	path, err := historyPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// PickedSet は選択済みの国を set として返す
func (h *History) PickedSet() map[string]bool {
	set := make(map[string]bool, len(h.Picked))
	for _, c := range h.Picked {
		set[c] = true
	}
	return set
}

// Add は国を履歴に追加する
func (h *History) Add(countries ...string) {
	h.Picked = append(h.Picked, countries...)
}

// Reset は履歴をクリアする
func (h *History) Reset() {
	h.Picked = []string{}
}
