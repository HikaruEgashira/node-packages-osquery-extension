# クイックスタート

このガイドに従って、osquery extensionを起動してください。

## ステップ1: ビルド

```bash
go build -o osquery-node-packages
```

## ステップ2: エクステンションを起動

### 開発時

```bash
osqueryi --allow_unsafe --extension ./osquery-node-packages
```

起動したら、以下のクエリでパッケージを確認できます：

```sql
-- パッケージマネージャーごとの数を確認
SELECT manager, COUNT(*) as count FROM node_packages GROUP BY manager;
```

期待される出力：
```
+---------+-------+
| manager | count |
+---------+-------+
| npm     | 1133  |
| bun     | 579   |
+---------+-------+
```

### プロダクション環境

プロダクション環境では、バイナリをroot所有にする必要があります：

```bash
# macOS
sudo chown root:wheel osquery-node-packages
osqueryi --extension ./osquery-node-packages

# Linux
sudo chown root:root osquery-node-packages
osqueryi --extension ./osquery-node-packages
```

## その他のクエリ例

```sql
-- エクステンションが登録されているか確認
SELECT * FROM osquery_extensions;

-- npmパッケージを10件表示
SELECT name, version FROM node_packages WHERE manager = 'npm' LIMIT 10;

-- bunパッケージを10件表示
SELECT name, version FROM node_packages WHERE manager = 'bun' LIMIT 10;

-- すべてのパッケージを表示
SELECT * FROM node_packages;
```

## トラブルシューティング

### エラー: "Extension binary has unsafe permissions"

開発時は`--allow_unsafe`フラグを使用してください：

```bash
osqueryi --allow_unsafe --extension ./osquery-node-packages
```

プロダクション環境では、バイナリの所有権を変更してください：

```bash
# macOS
sudo chown root:wheel osquery-node-packages

# Linux
sudo chown root:root osquery-node-packages
```

### エラー: "no such table: node_packages"

エクステンションが正しく起動していません。以下を確認してください：

```sql
SELECT * FROM osquery_extensions;
```

エクステンションが表示されない場合は、起動時のエラーメッセージを確認してください。

## 完全な実行例

```bash
$ osqueryi --allow_unsafe --extension ./osquery-node-packages
Using a virtual database. Need help, type '.help'
osquery> SELECT manager, COUNT(*) FROM node_packages GROUP BY manager;
+---------+-------+
| manager | count |
+---------+-------+
| npm     | 1133  |
| bun     | 579   |
+---------+-------+
osquery>
```

詳細は [USAGE.md](USAGE.md) を参照してください。
