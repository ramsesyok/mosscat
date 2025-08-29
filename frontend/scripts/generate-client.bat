@echo off
REM エラー時に処理を中断する
setlocal enabledelayedexpansion

REM openapi-typescript-codegen を実行
npx --yes openapi-typescript-codegen ^
  --input ..\internal\api\openapi.yaml ^
  --output src\api ^
  --client axios ^
  --useOptions ^
  --useUnionTypes

endlocal
