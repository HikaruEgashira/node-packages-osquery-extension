# osquery extensionの使い方

このツールはosquery extensionです。`osqueryi --extension`コマンドで起動します。

## セットアップ

### 1. ビルド

```bash
make build
# または
go build -o osquery-node-packages
```

## 使用方法

### 開発時

```bash
osqueryi --allow_unsafe --extension ./osquery-node-packages
```

`--allow_unsafe`フラグを使用すると、root所有でないバイナリでも起動できます。

### プロダクション環境

プロダクション環境では、セキュリティのためバイナリをroot所有にしてください：

```bash
# macOS
sudo chown root:wheel osquery-node-packages
osqueryi --extension ./osquery-node-packages

# Linux
sudo chown root:root osquery-node-packages
osqueryi --extension ./osquery-node-packages
```

## クエリ例

エクステンションが起動したら、以下のクエリを実行できます：

```sql
-- エクステンションが登録されているか確認
SELECT * FROM osquery_extensions;

-- すべてのパッケージを表示
SELECT * FROM node_packages;

-- npmパッケージのみ表示
SELECT name, version, cache_path
FROM node_packages
WHERE manager = 'npm'
LIMIT 10;

-- bunパッケージのみ表示
SELECT name, version, cache_path
FROM node_packages
WHERE manager = 'bun'
LIMIT 10;

-- パッケージマネージャーごとの数
SELECT manager, COUNT(*) as count
FROM node_packages
GROUP BY manager;
```

### クエリ例の出力

```
osquery> SELECT manager, COUNT(*) as count FROM node_packages GROUP BY manager;
+─────────+───────+
| manager | count |
+─────────+───────+
| npm     | 1133  |
| bun     | 579   |
+─────────+───────+
```

## テーブルスキーマ

`node_packages`テーブルは以下のカラムを提供します：

| カラム名 | 型 | 説明 |
|---------|-----|------|
| name | TEXT | パッケージ名 |
| version | TEXT | パッケージバージョン |
| manager | TEXT | パッケージマネージャー (npm, bun, pnpm, yarn, deno, jsr) |
| cache_path | TEXT | package.jsonへのパス |

## 検出されるキャッシュパス

### npm
- `~/.npm`
- `/usr/local/lib/node_modules`
- `/usr/lib/node_modules`

### bun
- `~/.bun/install/cache`
- `~/.bun/install/global`
- `~/.cache/.bun/install/cache`

### pnpm
- `~/.pnpm-store`
- `~/.local/share/pnpm/store`

### yarn
- `~/.yarn-cache`
- `~/.cache/yarn`

### deno
- `~/.cache/deno/npm`

### jsr
- `~/.cache/deno/deps/https/jsr.io`

## パフォーマンス

- 検出パッケージ数: 約1,700パッケージ（環境による）
  - npm: ~1,133パッケージ
  - bun: ~579パッケージ
- 初回スキャン時間: 約0.5秒
- ディレクトリ走査なし: 高速

## トラブルシューティング

### エラー: "Extension binary has unsafe permissions"

このエラーはosqueryのセキュリティチェックによるものです。

**開発時の解決方法:**
```bash
osqueryi --allow_unsafe --extension ./osquery-node-packages
```

**プロダクション環境での解決方法:**
```bash
# macOS
sudo chown root:wheel osquery-node-packages

# Linux
sudo chown root:root osquery-node-packages
```

### エラー: "no such table: node_packages"

エクステンションが正しく起動していません。以下を確認してください：

1. エクステンションが起動しているか
2. 正しいパスを使用しているか

エクステンションが登録されているか確認：

```sql
SELECT * FROM osquery_extensions;
```

### デバッグログの有効化

`--verbose`フラグを使用してデバッグログを有効化できます：

```bash
osqueryi --allow_unsafe --extension ./osquery-node-packages --verbose
```

## システム全体へのインストール

エクステンションをシステム全体で使用する場合：

```bash
# バイナリをインストール
sudo make install

# /etc/osquery/extensions.load に追加
echo "/usr/local/bin/osquery-node-packages" | sudo tee -a /etc/osquery/extensions.load

# osquerydを再起動
sudo systemctl restart osqueryd
```

## Tips

### ワンライナーでクエリ実行

```bash
echo "SELECT manager, COUNT(*) FROM node_packages GROUP BY manager;" | osqueryi --allow_unsafe --extension ./osquery-node-packages
```

### エクステンションを停止

Ctrl+Cを押すか：

```bash
pkill -f osquery-node-packages
```
