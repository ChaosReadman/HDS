# HDS

## ダミーサーバ
serverRoot以下にフォルダを作り、戻り値をretVal.jsonに記述することで、
任意のアクセスに対して戻り値を設定できる
Request時のパラメータ解析等はしていないので、
戻り値にパラメータ値を反映させるなどの機能は無い

## https接続対応
一応認証ファイルを使ってhttps接続できるようになている
オレオレ証明書なので、ブラウザからアクセスすると最初に信用しますか？などの警告が出る
curlでリクエストを投げる時は--insecureのオプションが必要

## 認証ファイルの作り方
認証ファイルはauthフォルダに入れてある
名前を合わせてKeyとCertを用意する

### 認証ファイルは以下の手順で作る

+ PowerShellで以下を実行

New-SelfSignedCertificate -DnsName PC-002 -FriendlyName "ftps-server" -CertStoreLocation "cert:\LocalMachine\My" -NotAfter (Get-Date).AddYears(10)


IISマネージャから、作成したファイルを選択してエクスポート  
test.pfxとして出力する  

次にOpenSSLをインストールし、コマンドプロンプトから以下を実行しpemを出力する  
openssl pkcs12 -in test.pfx -out test.pem  
次のコマンドでp12を出力する  
openssl pkcs12 -export -in test.pem -out test.p12

Goのサーバで使うには以下でusercert.pemとuserkey.pemを作る
openssl pkcs12 -in test.p12 -clcerts -nokeys -out usercert.pem  
openssl pkcs12 -in test.p12 -nocerts -out userkey.pem -nodes  

出来たファイルをauthに納めてある