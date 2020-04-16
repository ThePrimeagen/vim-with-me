function fizzbuzz(length: number, map: {[key: string]: number}) {
    for (let i = 1; i < length; ++i) {
        let str = "";
        for (const k in map) {
            const strToPrint = k;
            if (i % map[strToPrint] == 0) {
                str += strToPrint;
            }
        }

        if (str.length > 0) {
            console.log(str);
        }
        else {
            console.log(i);
        }
    }
}

fizzbuzz(136, {
    fizz: 3,
    buzz: 5,
    bazz: 7,
});
