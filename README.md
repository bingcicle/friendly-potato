# jsonl-tool
Merge or filter JSON Lines.

## Build & Run
```bash
go build -o jsonl-tool
./jsonl-tool --mode merge --in a.jsonl,b.jsonl > merged.jsonl
./jsonl-tool --mode filter --in merged.jsonl --field exchange --eq bybit
./jsonl-tool --mode filter --in merged.jsonl --field symbol --rex "BTC|ETH"
```
