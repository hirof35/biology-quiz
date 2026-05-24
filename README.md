# Go言語 × 生物学クイズ (Pure Go GUI Application)

Go言語（v1.16以上推奨）と、純粋なGo（Pure Go）だけでクロスプラットフォームなGUIを構築できるライブラリ `gioui.org` を使用した、中学生向けの生物学クイズアプリケーションです。

CGO（C Go）を完全に無効化した環境（`CGO_ENABLED=0`）で動作するため、GCCなどのC言語用コンパイラをインストールすることなく、軽量・高速なネイティブデスクトップアプリをビルド・実行できます。

---
<img width="758" height="675" alt="スクリーンショット 2026-05-25 044917" src="https://github.com/user-attachments/assets/9df3c29a-83d4-434a-b35f-a4be77ed0f1d" />

## 🚀 特徴

- **CGO不要 (Pure Go)**: 複雑な環境構築なしで、Windows、Mac、Linuxで即座に動作します。
- **データ駆動型**: クイズの問題（全30問）はすべて外部の `quiz.json` から動的に読み込まれます。問題の追加や修正もJSONファイルを書き換えるだけで容易に行えます。
- **モダンなUI**: 状態変化に応じて即座に描画を更新する、即時モード（Immediate Mode）GUIアーキテクチャを採用しています。

---

## 📂 ディレクトリ構成

プロジェクトフォルダ内は、以下の配置にしてください。

```text
biology-quiz/
├── main.go      # GUIアプリケーション本体のソースコード
├── quiz.json    # 中学生理科（生物分野）を網羅した30問のクイズデータ
└── README.md    # 本ドキュメント
🛠️ 開発環境のセットアップと実行方法
Windows (PowerShell) 環境での実行手順です。

1. モジュールの初期化
プロジェクトのルートディレクトリで以下のコマンドを実行し、Goのモジュールシステムを初期化します。

PowerShell
go mod init biology-quiz
2. 依存パッケージ（Gio UI）の取得
Gio UIおよびその関連サブモジュール（シェーダーやフォント等）を一括で強制ダウンロードします。

PowerShell
go get -u gioui.org/...
3. 依存関係のクリーンアップ
go.mod と go.sum の整合性を整えます。

PowerShell
go mod tidy
4. アプリケーションの実行
CGOを明示的に無効化して、プログラムを起動します。

PowerShell
$env:CGO_ENABLED=0; go run main.go
📑 クイズデータの仕様 (quiz.json)
問題データを追加・変更する場合は、quiz.json を以下の構造に従って編集してください。

questionNumber: 問題番号 (整数)

question: 問題文 (文字列)

answerOptions: 4つの選択肢の配列。

text: 選択肢のテキスト

isCorrect: 正解なら true、不正解なら false (各問題に必ず1つだけ true を配置)

🔧 技術スタック・主要な仕様
Language: Go

GUI Framework: gioui.org (v0.10.0 仕様準拠)

Architecture: Immediate Mode GUI (即時モード)

JSON Parser: encoding/json (標準パッケージ)
