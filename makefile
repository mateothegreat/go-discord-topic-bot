include .env
export

.PHONY: docs

docs:
	ffmpeg -y -i docs/demo.mp4 -vf fps=30,scale=1080:-1:flags=lanczos,palettegen palette.png -y
	ffmpeg -y -i docs/demo.mp4 -i palette.png -filter_complex "fps=30,scale=1080:-1:flags=lanczos[x];[x][1:v]paletteuse" docs/demo.gif -y
	rm palette.png

db/generate:
	cd prisma && go run github.com/steebchen/prisma-client-go generate

db/migrate:
	cd prisma && go run github.com/steebchen/prisma-client-go db push
