import TCP from "../tcp";

test("TCP#env", function() {
    const buffer = TCP.createEnv("testing");
    const expected = Buffer.alloc(7 + 2);
    expected.writeUInt16BE("testing".length, 0);
    Buffer.from("testing").copy(expected, 2);

    expect(buffer).toEqual(expected);
});

