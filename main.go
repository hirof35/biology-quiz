package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// JSONデータ構造
type Option struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

type Question struct {
	QuestionNumber int      `json:"questionNumber"`
	Question       string   `json:"question"`
	AnswerOptions  []Option `json:"answerOptions"`
}

type QuizData struct {
	Questions []Question `json:"questions"`
}

// クイズの状態管理
type QuizSystem struct {
	data          QuizData
	currentIndex  int
	score         int
	answered      bool
	lastResult    string
	optionButtons []widget.Clickable
	nextButton    widget.Clickable
}

func main() {
	// 1. JSONの読み込み
	jsonFile, err := os.Open("quiz.json")
	if err != nil {
		log.Fatalf("JSONファイルが開けません: %v", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var quizData QuizData
	if err := json.Unmarshal(byteValue, &quizData); err != nil {
		log.Fatalf("JSONの解析に失敗しました: %v", err)
	}

	// システム状態の初期化
	sys := &QuizSystem{
		data:          quizData,
		currentIndex:  0,
		score:         0,
		answered:      false,
		optionButtons: make([]widget.Clickable, 4),
	}

	// 2. Gio アプリケーションの起動
	go func() {
		var w app.Window
		w.Option(app.Title("Go言語×生物学クイズ (Pure Go版)"))
		w.Option(app.Size(unit.Dp(600), unit.Dp(500)))

		if err := sys.run(&w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// メインの描画ループ
func (sys *QuizSystem) run(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	for {
		e := w.Event()
		switch ev := e.(type) {
		case app.DestroyEvent:
			return ev.Err
		case app.FrameEvent:
			// コンテキストの生成
			gtx := app.NewContext(&ops, ev)

			// 状態更新
			sys.update(gtx, w)

			// 画面レイアウトの構築
			layout.Flex{
				Axis:      layout.Vertical,
				Alignment: layout.Middle,
			}.Layout(gtx,
				// 上部：進行状況
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					var progressText string
					if sys.currentIndex < len(sys.data.Questions) {
						progressText = fmt.Sprintf("第 %d / %d 問", sys.currentIndex+1, len(sys.data.Questions))
					} else {
						progressText = "🎉 クイズ終了！"
					}
					return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						label := material.H5(th, progressText)
						label.Alignment = text.Middle
						return label.Layout(gtx)
					})
				}),

				// 中央：問題文
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					var qText string
					if sys.currentIndex < len(sys.data.Questions) {
						qText = sys.data.Questions[sys.currentIndex].Question
					} else {
						qText = fmt.Sprintf("すべての問題が終了しました！\n\nあなたのスコア: %d / %d 点", sys.score, len(sys.data.Questions))
					}
					return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						label := material.Body1(th, qText)
						return label.Layout(gtx)
					})
				}),

				// 下部：選択肢ボタン群
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if sys.currentIndex >= len(sys.data.Questions) {
							return layout.Dimensions{}
						}

						currentQ := sys.data.Questions[sys.currentIndex]

						var kids []layout.FlexChild
						for i, option := range currentQ.AnswerOptions {
							idx := i
							opt := option
							kids = append(kids, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									btn := material.Button(th, &sys.optionButtons[idx], opt.Text)
									
									// 最新仕様: 解答済みの場合はクリックイベントを無効化（再描画コマンドを明示的に送る）
									if sys.answered {
										gtx.Execute(op.InvalidateCmd{})
									}
									return btn.Layout(gtx)
								})
							}))
						}

						// 解答済みなら結果テキストと「次へ」ボタンを追加
						if sys.answered {
							kids = append(kids, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: unit.Dp(15), Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									lbl := material.H6(th, sys.lastResult)
									lbl.Alignment = text.Middle
									return lbl.Layout(gtx)
								})
							}))
							kids = append(kids, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return material.Button(th, &sys.nextButton, "次の問題へ").Layout(gtx)
							}))
						}

						return layout.Flex{Axis: layout.Vertical}.Layout(gtx, kids...)
					})
				}),
			)

			ev.Frame(&ops)
		}
	}
}

// 内部状態の更新メソッド
func (sys *QuizSystem) update(gtx layout.Context, w *app.Window) {
	if sys.currentIndex >= len(sys.data.Questions) {
		return
	}

	// 1. 選択肢ボタンのクリック検知
	if !sys.answered {
		currentQ := sys.data.Questions[sys.currentIndex]
		for i := range sys.optionButtons {
			if sys.optionButtons[i].Clicked(gtx) {
				sys.answered = true
				if currentQ.AnswerOptions[i].IsCorrect {
					sys.score++
					sys.lastResult = "⭕ 正解！"
				} else {
					sys.lastResult = "❌ 不正解..."
				}
				w.Invalidate()
			}
		}
	}

	// 2. 「次の問題へ」ボタンのクリック検知
	if sys.answered && sys.nextButton.Clicked(gtx) {
		sys.currentIndex++
		sys.answered = false
		sys.lastResult = ""
		sys.optionButtons = make([]widget.Clickable, 4) // 状態リセット
		w.Invalidate()
	}
}