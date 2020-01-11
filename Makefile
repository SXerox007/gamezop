.PHONY: all app rest clean

app:
	go run server/myapplication.go

rest:
	go run server/rest-app/rest-app.go