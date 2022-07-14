# Dinkur the task time tracking utility.
# <https://github.com/dinkur/dinkur>
#
# SPDX-FileCopyrightText: 2021 Kalle Fagerberg
# SPDX-License-Identifier: CC0-1.0

ifeq ($(OS),Windows_NT)
dinkur-gnome-raujonas-executor.exe:
else
dinkur-gnome-raujonas-executor:
endif
	go build

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -rfv ./dinkur-statusline.exe ./dinkur-statusline

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: deps
deps:
	go get
