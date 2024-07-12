Dieses Programm synchronisiert alle Artikel auf dem SAGE für den Labeldrucker

Zuerst Go runterladen und installieren
https://go.dev/dl/

Danach das Repo runterladen
https://github.com/computerextra/P-touch_Labelsync-go

In den Ordner gehen, wo das Zeug runtergeladen wurde:
Terminal in dem Ordner öffnen:

```pwsh
go get .
```

Warten bis fertig.

Danach
`.env` Datei anlegen und nach dem example ausfüllen.

```
go build
```

In dem Ordner ist nun eine `labelsync.exe` Datei.

Wenn diese Ausgeführt wird, passiert alles von alleine.

Die oberen Schritte müssen nur einmalig durchgeführt werden, danach kann direkt die `labelsync.exe` geöffnet werden.
