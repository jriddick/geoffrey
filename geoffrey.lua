-- Add a bot
geoffrey.add("geoffrey", {
    Registration = {
        Nick = "geoffrey",
        User = "geoffrey",
        Name = "geoffrey",
    },
    Authentication = {
        Username = "geoffrey",
        Password = "...",
    },
    Limits = {
        Reconnect = 5,
        Timeout = 5,
    },
    Connection = {
        Host = "localhost",
        Port = 6667,
        Secure = false,
        InsecureSkipVerify = false,
        Timeout = 30
    },
    Channels = {
        "#geoffrey"
    },
    Plugins = {
        "echo",
        "registration",
        "join",
        "ping"
    },
})