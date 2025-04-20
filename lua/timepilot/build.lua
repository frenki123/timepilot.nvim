M = {}

function M.build()
  coroutine.yield("Installing timepilot CLI...")
  if vim.fn.executable("go") == 1 then
    vim.system({
      "go",
      "install",
      "github.com/frenki123/timepilot.nvim/timepilot@latest",
    }, {
      stdout = false,
      stderr = false,
    }, function(obj)
      if obj.code == 0 then
        print("✅ timepilot CLI installed successfully.")
      else
        print("❌ Failed to install timepilot CLI. Exit code:", obj.code)
      end
    end)
  else
    print("⚠ Go is not installed. Cannot install timepilot CLI.")
  end
end

return M
