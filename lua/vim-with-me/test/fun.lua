local test_count = {}
for idx = 1, #VWM_giveaway_entries do
    if test_count[VWM_giveaway_entries[idx]] == nil then
        test_count[VWM_giveaway_entries[idx]] = 0
    end

    test_count[VWM_giveaway_entries[idx]] = test_count[VWM_giveaway_entries[idx]] + 1
end

print("Test", vim.inspect(test_count))

