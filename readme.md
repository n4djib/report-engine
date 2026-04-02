== Initialize GO ==
go mod init github.com/n4djib/reporter
    initialize mod file
go mod tidy
    tidy up mod

== GIT
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

git init


== Folder Structure ==
/cmd : contains the apps
/pkg : is for code you are OK with other projects importing.
/internal : is for code that must stay private inside your repo
/rest : is for rest api tests
/tmp : air build output path

== Air ==
go install github.com/air-verse/air@latest

air init

after we configure AIR

== Run using Makefile.cmd
Makefile.cmd            display help
Makefile.cmd dc         run Services in docker-compose
Makefile.cmd sd         run Server with AIR
Makefile.cmd rd         run Remote with AIR



== nk (nkey key pair generator)
go install github.com/nats-io/nkeys/nk@latest

nk -gen user -pubout

we can generate the pair using "github.com/nats-io/nkeys"




