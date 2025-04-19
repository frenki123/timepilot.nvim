local M ={}

local function install_timepilot()
    if M.config.autoinstall then
        return os.execute("go install github.com/frenki123/timepilot.nvim/timepilot@latest")
    end
    return false
end

function M.build()
    local required_version = "0.0.1"
    local timepilot = M.config.timepilot_path
    local handle = io.popen(timepilot .. "version")
    if not handle then
        if not install_timepilot() then
            print("Timepilot not installed")
            return false
        end
        return
    end
    local result = handle:read("*a")
    handle:close()
    local version = result:match("v(.+)")
    if version == required_version then
        return
    end
    if not install_timepilot() then
        print("Required Timepilot version is" .. required_version .. "installed" .. version)
    end
end


M.config = {
    timepilot_path = "timepilot",
    autoinstall = true
}


function M.set_config(user_config)
    for option, value in pairs(user_config) do
        M.config[option] = value
    end
end

function M.check_config()
end

return M
