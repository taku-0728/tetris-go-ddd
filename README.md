# テトリスゲーム（Go言語 + DDD実装）

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

ドメイン駆動設計（DDD）アーキテクチャを採用したテトリスゲームのGo言語実装です。クリーンアーキテクチャの実践例として、教育的価値の高いプロジェクトです。

## 🎮 機能

### ゲーム機能
- **完全なテトリス実装**: 7種類のテトロミノ（I, O, T, S, Z, J, L）
- **標準ゲームボード**: 10×20のプレイフィールド
- **完全な操作系**: 移動、回転、落下、一気落下
- **ライン消去**: 完成したラインの自動消去とスコア計算
- **レベルシステム**: プレイ進行に応じた難易度調整
- **ゲームオーバー判定**: 適切な終了条件

### システム機能
- **リアルタイム処理**: 60FPS ゲームループ
- **非同期入力処理**: ゴルーチンベースの応答性の高い入力
- **一時停止/再開**: ゲーム中断機能
- **リスタート**: ゲーム再開機能
- **包括的エラーハンドリング**: 全層での堅牢なエラー処理

## 🏗️ アーキテクチャ

### DDD層構造
```
tetris/
├── presentation/           # プレゼンテーション層
│   └── main.go            # ゲームループとメイン関数
├── application/           # アプリケーション層
│   └── game_controller.go # ゲーム制御ロジック
├── domain/               # ドメイン層
│   ├── model/           # ドメインモデル
│   │   ├── point.go     # 座標値オブジェクト
│   │   ├── board.go     # ゲームボード
│   │   └── tetromino.go # テトロミノ
│   └── service/         # ドメインサービス
│       └── game_service.go # ゲームコアロジック
└── infrastructure/      # インフラストラクチャ層
    ├── console/         # コンソール表示
    │   └── display.go
    └── input/          # 入力処理
        └── keyboard.go
```

### 設計原則
- **単一責任原則**: 各層・各クラスが明確な責任を持つ
- **依存性逆転**: 上位層が下位層の詳細に依存しない
- **ドメイン中心**: ビジネスロジックをドメイン層に集約
- **テスト駆動**: 全層で包括的なテスト実装

## 🚀 クイックスタート

### 前提条件
- Go 1.23 以上
- Unix系OS（macOS/Linux）またはWindows

### インストール・実行
```bash
# リポジトリのクローン
git clone <repository-url>
cd claude_study

# 依存関係の確認
go mod tidy

# ゲーム実行
go run presentation/main.go
```

### ビルド
```bash
# 実行可能ファイルの作成
go build -o tetris presentation/main.go

# 実行
./tetris
```

## 🎯 操作方法

| キー | 動作 |
|------|------|
| `A` / `D` | 左右移動 |
| `S` | 下移動 |
| `W` | 回転 |
| `Space` | 一気に落下 |
| `P` | 一時停止/再開 |
| `Q` | 終了 |
| `R` | リスタート |

## 🧪 テスト

### テスト実行
```bash
# 全テスト実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 特定パッケージのテスト
go test ./domain/model
go test ./domain/service
go test ./application
go test ./infrastructure/input
```

### テスト戦略
- **テーブル駆動テスト**: 全テストでテーブル駆動方式を採用
- **包括的カバレッジ**: 正常ケース・エラーケース・エッジケースを網羅
- **層別テスト**: 各アーキテクチャ層で独立したテスト

### テストカバレッジ
```
application/        76.9% coverage
domain/model/       83.5% coverage  
domain/service/     68.4% coverage
infrastructure/input/ 57.1% coverage
```

## 📁 プロジェクト構成

### 主要ファイル
```
claude_study/
├── go.mod                           # Go モジュール定義
├── CLAUDE.md                        # 開発者向け情報
├── README.md                        # プロジェクト説明（本ファイル）
├── presentation/main.go             # メインゲームループ（319行）
├── application/
│   ├── game_controller.go           # ゲーム制御（172行）
│   └── game_controller_test.go      # テスト（467行）
├── domain/
│   ├── model/
│   │   ├── point.go                 # 座標（9行）
│   │   ├── point_test.go            # テスト（33行）
│   │   ├── board.go                 # ボード（148行）
│   │   ├── board_test.go            # テスト（481行）
│   │   ├── tetromino.go             # テトロミノ（250行）
│   │   └── tetromino_test.go        # テスト（339行）
│   └── service/
│       ├── game_service.go          # ゲームロジック（234行）
│       └── game_service_test.go     # テスト（390行）
└── infrastructure/
    ├── console/
    │   └── display.go               # 表示処理（138行）
    └── input/
        ├── keyboard.go              # 入力処理（136行）
        └── keyboard_test.go         # テスト（395行）
```

### 技術的詳細

#### ドメインモデル
- **Point**: 座標を表す値オブジェクト
- **Tetromino**: テトロミノの形状・位置・回転状態を管理
- **Board**: ゲームボードの状態とライン消去ロジック

#### ドメインサービス  
- **GameService**: テトリスのコアビジネスロジック
  - ピース移動・回転・落下処理
  - ライン消去とスコア計算
  - ゲームオーバー判定

#### アプリケーション層
- **GameController**: ゲーム制御とタイマー管理
  - 入力処理の振り分け
  - ゲーム状態の管理
  - 落下タイミングの制御

#### インフラストラクチャ層
- **Display**: コンソール出力とゲーム画面描画
- **KeyboardInput**: 非同期キーボード入力処理

## 🔧 開発

### 開発環境セットアップ
```bash
# プロジェクトルートで
go mod download

# テスト環境確認
go test ./... -v
```

### コードフォーマット
```bash
# フォーマット実行
go fmt ./...

# インポート整理
goimports -w .
```

### 静的解析
```bash
# vet実行
go vet ./...

# golangci-lint（オプション）
golangci-lint run
```

## 📚 学習ポイント

このプロジェクトは以下の概念の学習に最適です：

### アーキテクチャパターン
- **ドメイン駆動設計（DDD）**
- **クリーンアーキテクチャ**
- **層化アーキテクチャ**

### Go言語の機能
- **ゴルーチンとチャネル**: 非同期入力処理
- **インターフェース**: 依存性の抽象化
- **エラーハンドリング**: 適切なエラー処理パターン
- **テスト**: テーブル駆動テストの実装

### 設計原則
- **SOLID原則**の実践
- **関心の分離**
- **依存性注入**
- **テスト可能な設計**

## 🤝 コントリビューション

プルリクエストや課題報告は歓迎します。大きな変更を行う場合は、まずIssueで議論をお願いします。

### 開発ガイドライン
1. テスト駆動開発を心がける
2. アーキテクチャの原則を遵守する
3. 適切なエラーハンドリングを実装する
4. 日本語でのコメント・ドキュメント作成

## 📄 ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は[LICENSE](LICENSE)ファイルを参照してください。

## 🎯 今後の拡張予定

- [ ] ネットワーク対戦機能
- [ ] AIプレイヤー
- [ ] ハイスコア保存機能
- [ ] グラフィカルUI（Webベース）
- [ ] モバイル対応

---
