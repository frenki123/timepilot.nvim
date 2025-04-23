M = {}

local client = require("timepilot.client")
local config = require("timepilot.config").config

local uv = vim.uv
local project
local idle = false
local idle_timeout = config.timeout * 60 * 1000
local idle_timer

local function match_filetype(ft)
  for _, value in pairs(config.disabled_filetypes) do
    local matcher = "^" .. value .. (value:sub(-1) == "*" and "" or "$")
    if ft:match(matcher) then
      return true
    end
  end

  return false
end

local function should_ignore()
  return vim.tbl_contains(config.disabled_filetypes, vim.bo.ft)
    or match_filetype(vim.bo.ft)
    or vim.api.nvim_get_option_value("buftype", { buf = 0 }) == "terminal"
    or vim.fn.reg_executing() ~= ""
    or vim.fn.reg_recording() ~= ""
end

local function run_cmd(cmd, cmd_args, cwd, callback)
  local stdout = uv.new_pipe(false)
  local stderr = uv.new_pipe(false)
  local handle = uv.spawn(cmd, {
    args = cmd_args,
    cwd = cwd,
    stdio = { nil, stdout, stderr },
  }, function(_, _)
    stdout:close()
    stderr:close()
  end)
  if not handle then
    return
  end
  handle:close()
  local result = {}
  stdout:read_start(function(err, data)
    if err then
      return
    end
    if data then
      table.insert(result, data)
    else
      vim.schedule(function()
        callback(table.concat(result))
      end)
    end
  end)
end

local function on_idle()
  if idle then
    return
  end -- already idle

  idle = true
  local filename = vim.api.nvim_buf_get_name(0)
  client.send_event("session/leave", { project = project })
  if not should_ignore() then
    client.send_event("buffer/leave", { project = project, filename = filename })
  end
end

local function on_activity()
  if idle then
    -- Back from the dead
    idle = false
    local filetype = vim.bo.filetype
    local filename = vim.api.nvim_buf_get_name(0)
    client.send_event("session/enter", { project = project })
    if not should_ignore() then
      client.send_event("buffer/enter", { project = project, filename = filename, filetype = filetype })
    end
  end

  -- Reset timer
  if idle_timer then
    idle_timer:stop()
    idle_timer:close()
  end
  idle_timer = uv.new_timer()
  idle_timer:start(idle_timeout, 0, vim.schedule_wrap(on_idle))
end

function M.autocmd()
  local augroup = vim.api.nvim_create_augroup("Timepilot", { clear = true })
  --Start client
  vim.api.nvim_create_autocmd("VimEnter", {
    group = augroup,
    once = true,
    callback = function()
      client.start()
      on_activity()
      vim.on_key(function()
        on_activity()
      end, vim.api.nvim_create_namespace("timepilot_idle"))
    end,
  })

  vim.api.nvim_create_autocmd("VimEnter", {
    group = augroup,
    callback = function()
      run_cmd("git", { "rev-parse", "--show-toplevel" }, vim.fn.getcwd(), function(output)
        project = output or vim.fn.getcwd()
        project:match("^%s*(.-)%s*$")
        vim.defer_fn(function()
          client.send_event("session/enter", { project = project })
        end, 250)
      end)
    end,
  })

  vim.api.nvim_create_autocmd("VimLeavePre", {
    group = augroup,
    callback = function()
      client.send_event("session/leave", { project = project })
    end,
  })

  vim.api.nvim_create_autocmd("BufWinEnter", {
    group = augroup,
    callback = function()
      if should_ignore() then
        return
      end
      local filetype = vim.bo.filetype
      local filename = vim.api.nvim_buf_get_name(0)
      local buftype = vim.bo.buftype
      if filename == "" or buftype ~= "" then
        return
      end
      client.send_event("buffer/enter", { project = project, filename = filename, filetype = filetype })
    end,
  })

  vim.api.nvim_create_autocmd("BufWinLeave", {
    group = augroup,
    callback = function()
      if should_ignore() then
        return
      end
      local filename = vim.api.nvim_buf_get_name(0)
      local buftype = vim.bo.buftype
      if filename == "" or buftype ~= "" then
        return
      end
      client.send_event("buffer/leave", { project = project, filename = filename })
    end,
  })
end

function M.set_keys()
  vim.api.nvim_create_user_command("GetProjectTime", function()
    client.send_event("data/project", { project_name = project })
  end, {})
  vim.api.nvim_create_user_command("GetFileTime", function()
    client.send_event("data/file", { project_name = project })
  end, {})
end

return M
