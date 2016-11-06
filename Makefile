NAME = geoffrey
VERSION = 0.1.0
EXTRA = const.lua config.lua geoffrey.lua plugins/*
WINDOWS := win64 win32
LINUX := linux32 linux64
DARWIN := darwin
ARCH := $(WINDOWS) $(LINUX) $(DARWIN)
ZIPS = $(NAME)-$(VERSION)-$(ARCH).zip

$(NAME)-win64.exe:
	GOOS=windows GOARCH=amd64 go build -o $(NAME)-win64.exe geoffrey.go

$(NAME)-win32.exe:
	GOOS=windows GOARCH=386 go build -o $(NAME)-win32.exe geoffrey.go

$(NAME)-linux32:
	GOOS=linux GOARCH=386 go build -o $(NAME)-linux32 geoffrey.go

$(NAME)-linux64:
	GOOS=linux GOARCH=amd64 go build -o $(NAME)-linux64 geoffrey.go

$(NAME)-darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(NAME)-darwin geoffrey.go

$(NAME)-$(VERSION)-win32.zip: $(addprefix $(NAME),-win32.exe)
	zip $(NAME)-$(VERSION)-win32.zip geoffrey-win32.exe $(EXTRA)

$(NAME)-$(VERSION)-win64.zip: $(addprefix $(NAME),-win64.exe)
	zip $(NAME)-$(VERSION)-win64.zip geoffrey-win64.exe $(EXTRA)

$(NAME)-$(VERSION)-linux32.zip: $(addprefix $(NAME),-linux32)
	zip $(NAME)-$(VERSION)-linux32.zip geoffrey-linux32 $(EXTRA)

$(NAME)-$(VERSION)-linux64.zip: $(addprefix $(NAME),-linux64)
	zip $(NAME)-$(VERSION)-linux64.zip geoffrey-linux64 $(EXTRA)

$(NAME)-$(VERSION)-darwin.zip: $(addprefix $(NAME),-darwin)
	zip $(NAME)-$(VERSION)-darwin.zip geoffrey-darwin $(EXTRA)

zip: $(addprefix $(NAME)-$(VERSION)-,$(addsuffix .zip,$(ARCH)))
