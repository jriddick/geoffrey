-- Add a bot
geoffrey.add("geoffrey", {
    Config = {
        Hostname = "localhost",
        Port = 6667,
        Secure = false,
        InsecureSkipVerify = false,
        Nick = "geoffrey",
        User = "geoffrey",
        Name = "geoffrey",
        ReconnectLimit = 2,
        Timeout = 30,
        TImeoutLimit = 5,
        Channels = {
            "#geoffrey"
        },
    },
    Plugins = {
        "echo"
    }
})