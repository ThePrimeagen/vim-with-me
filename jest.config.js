module.exports = function() {
    return {
        rootDir: "src",
        testRegex: "/__tests__/.*",
        preset: "ts-jest",
        transform: {"\\.ts$": ["ts-jest"]},
        moduleNameMapper: {
            "~/(.*)": "<rootDir>/$1",
        }
    };
};

