#!/usr/bin/env node

const log = require("./logger");

const names = [
    "foo",
    "bar",
    "baz",
]

const towerPoints = !!process.argv[2] ? [
    "t:0:0",
    "t:1:1",
    "t:2:4",
    "t:4:8",
    "t:6:9",
    "t:6:9",
] : [
    "t:0:0",
    "t:1:1",
    "t:2:4",
    "t:4:8",
]

/**
 * @param {any[]} list
 */
function rand(list) {
    return list[Math.floor(Math.random() * list.length)]
}

/**
 * @param {any[]} list
 */
function walk(list) {
    let idx = 0;
    return function() {
        idx++;
        return list[idx % list.length]
    }
}

const msgs = walk(towerPoints)
setInterval(function() {
    log(`message:${rand(names)}:${msgs()}`);
}, 200)


