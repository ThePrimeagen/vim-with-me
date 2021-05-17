jest.mock("../../now", () => ({
    default: jest.fn()
}));
import SystemCommand from "../index";
import accumulateBus from "~/utils/bus";
import now from "~/now";

const mockNow = now as jest.Mock;

jest.useFakeTimers();
describe("SystemCommand", function() {
    it("happy path, one command, emits on and off strings", function() {
        mockNow.mockImplementation(() => 23);

        const flushLogs = accumulateBus();
        const command = new SystemCommand("on", "off", 5000);

        command.add();

        expect(flushLogs()).toEqual([[
            "system-command", "on"
        ]]);

        mockNow.mockImplementation(() => 5000);
        jest.advanceTimersByTime(1000);

        expect(flushLogs()).toEqual([[
            "system-command", "on"
        ]]);

        mockNow.mockImplementation(() => 5024);
        jest.advanceTimersByTime(100);
        expect(flushLogs()).toEqual([[
            "system-command", "on"
        ], [
            "system-command", "off"
        ]]);
    });

    it("multiple commands", function() {
        mockNow.mockImplementation(() => 23);

        const flushLogs = accumulateBus();
        const command = new SystemCommand("on", "off", 5000);

        command.add();

        expect(flushLogs()).toEqual([[
            "system-command", "on"
        ]]);

        mockNow.mockImplementation(() => 5000);
        jest.advanceTimersByTime(1000);

        expect(flushLogs()).toEqual([[
            "system-command", "on"
        ]]);

        command.add();
        command.add();

        mockNow.mockImplementation(() => 5024);
        jest.advanceTimersByTime(100);
        expect(flushLogs()).toEqual([[
            "system-command", "on"
        ]]);

        mockNow.mockImplementation(() => 15100);
        jest.advanceTimersByTime(100);
        expect(flushLogs()).toEqual([[
            "system-command", "on"
        ], [
            "system-command", "off"
        ]]);
    });
});

