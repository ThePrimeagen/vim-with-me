const mockMath = Object.create(global.Math);
mockMath.random = () => 0.5;
global.Math = mockMath;

let ons: {[key: string]: ((...args: any[]) => void)[]} = {};
jest.mock("../../message-bus", () => ({ default: {
    emit: jest.fn(),
    on: (str: string, callback: (...args: any[]) => void) => {
        if (!ons[str]) {
            ons[str] = []
        }
        ons[str].push(callback);
    },
}}));

import bus from "../../message-bus";

import ProgramWithMe from "../index";

jest.useFakeTimers();
describe("ProgramWithMe", function() {

    beforeEach(function() {
        // @ts-ignore
        bus.emit.mockReset();
        ons = {};
    });

    it("should be able to add a programmer", function() {
        const pwm = new ProgramWithMe();

        ons["quirk-message"][0]({
            rewardName: "ProgrammWithMeEnter",
            username: "foo-bar",
        });

        expect(pwm.programmers).toEqual([
            "foo-bar",
        ]);

        ons["quirk-message"][0]({
            rewardName: "ProgrammWithMeEnter",
            username: "foo-bar",
        });

        expect(pwm.programmers).toEqual([
            "foo-bar",
        ]);

        ons["quirk-message"][0]({
            rewardName: "ProgrammWithMeEnter",
            username: "bar-foo",
        });

        expect(pwm.programmers).toEqual([
            "foo-bar",
            "bar-foo",
        ]);
    });

    it("should be able to add a programmer", function() {
        const pwm = new ProgramWithMe();

        ons["quirk-message"][0]({
            rewardName: "ProgrammWithMeEnter",
            username: "foo-bar",
        });

        ons["quirk-message"][0]({
            rewardName: "ProgrammWithMeEnter",
            username: "bar-foo",
        });

        expect(pwm.programmers).toEqual([
            "foo-bar",
            "bar-foo",
        ]);

        mockMath.random = () => 0.1;
        pwm.enableProgramWithMe();
        expect(pwm.currentFamousPerson).toEqual("foo-bar");
        //@ts-ignore
        expect(bus.emit).toHaveBeenCalled();

        mockMath.random = () => 0.7;
        jest.advanceTimersByTime(61000);
        expect(pwm.currentFamousPerson).toEqual("bar-foo");
        //@ts-ignore
        expect(bus.emit).toHaveBeenCalledTimes(2);

        mockMath.random = () => 0.1;
        pwm.validateFunction({
            rewardName: "VimInsert",
            username: "bar-foo",
            userInput: "",
            cost: 69,
        });
        expect(pwm.currentFamousPerson).toEqual("foo-bar");
        expect(bus.emit).toHaveBeenCalledTimes(3);
    });
});
