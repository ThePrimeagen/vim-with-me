import validate from "../index";

describe("vim-commands", function() {
    it("should dd", function() {
        expect(validate({
            rewardName: "Vim Command PWM pwm",
            userInput: "dd",
            cost: 69420,
            username: "theprimeagen",
        })).toEqual({
            success: true
        });
    });
});
