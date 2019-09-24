if (!String.prototype.startsWith) {
    String.prototype.startsWith = function (searchString, position) {
        "use strict";
        var str = this,
            strLen = str.length,
            seaLen = searchString.length,
            pos = position || 0,
            i;

        if (seaLen + pos > strLen) {
            return false;
        }

        for (i = 0; i < seaLen; i++) {
            if (str.charCodeAt(pos + i) !== searchString.charCodeAt(i)) {
                return false;
            }
        }
        return true;
    };
}
