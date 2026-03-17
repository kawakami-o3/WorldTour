package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	n := flag.Int("n", 1, "選択する国の数")
	list := flag.Bool("list", false, "選択済みの国一覧を表示")
	reset := flag.Bool("reset", false, "選択履歴をリセット")
	remaining := flag.Bool("remaining", false, "未選択の国数を表示")
	flag.Parse()

	history, err := LoadHistory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "履歴の読み込みに失敗しました: %v\n", err)
		os.Exit(1)
	}

	switch {
	case *list:
		cmdList(history)
	case *reset:
		cmdReset(history)
	case *remaining:
		cmdRemaining(history)
	default:
		cmdPick(history, *n)
	}
}

func cmdList(h *History) {
	if len(h.Picked) == 0 {
		fmt.Println("まだ国が選択されていません。")
		return
	}
	fmt.Printf("選択済みの国 (%d カ国):\n", len(h.Picked))
	for i, c := range h.Picked {
		fmt.Printf("  %3d. %s\n", i+1, c)
	}
}

func cmdReset(h *History) {
	h.Reset()
	if err := h.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "履歴のリセットに失敗しました: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("履歴をリセットしました。")
}

func cmdRemaining(h *History) {
	picked := h.PickedSet()
	count := 0
	for _, c := range Countries {
		if !picked[c] {
			count++
		}
	}
	fmt.Printf("未選択: %d カ国 / 全 %d カ国\n", count, len(Countries))
}

func cmdPick(h *History, n int) {
	picked := h.PickedSet()

	var candidates []string
	for _, c := range Countries {
		if !picked[c] {
			candidates = append(candidates, c)
		}
	}

	if len(candidates) == 0 {
		fmt.Println("すべての国を選び終えました！おめでとうございます！")
		return
	}

	if n > len(candidates) {
		fmt.Printf("残り %d カ国しかありません。すべて選択します。\n", len(candidates))
		n = len(candidates)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	selected := candidates[:n]

	fmt.Println("選ばれた国:")
	for _, c := range selected {
		fmt.Printf("  🌍 %s\n", c)
	}

	h.Add(selected...)
	if err := h.Save(); err != nil {
		fmt.Fprintf(os.Stderr, "履歴の保存に失敗しました: %v\n", err)
		os.Exit(1)
	}

	remainingCount := len(candidates) - n
	fmt.Printf("\n残り %d カ国\n", remainingCount)
}
