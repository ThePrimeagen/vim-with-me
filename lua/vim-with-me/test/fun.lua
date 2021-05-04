local test_count = {}
for idx = 1, #giveaway_entries do
    if test_count[giveaway_entries[idx]] == nil then
        test_count[giveaway_entries[idx]] = 0
    end

    test_count[giveaway_entries[idx]] = test_count[giveaway_entries[idx]] + 1
end

print("Test", vim.inspect(test_count))

