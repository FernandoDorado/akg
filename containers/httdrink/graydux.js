/* Title: Graydux
 * Author: Adam Gray
 * Created: 2018-08-09
 * Updated: 2018-09-30
 * Description: State management for tiny webapps
 */

function Graydux() {
    // List of functions which define how state changes
    this.reducers = [];

    // State of web application
    this.store = {};

    // By default, onReducerFail sends the error to the console
    this.onReducerFail = function (error) { console.log(error); }
}

Graydux.prototype.addReducer = function (name, keyList, method) {
    this.reducers = this.reducers.concat([{
        name: name,         // Name of the reducer
        keyList: keyList,   // Part of store which reducer manages
        method: method,     // The reducer function
        subscribers: []     // Functions which subscribe to this reducer's state change
    }]);
};

Graydux.prototype.subscribe = function (name, method) {
    var i;
    var reducerCount = this.reducers.length;
    // For each reducer ...
    for (i = 0; i < reducerCount; i ++) {
        // ... if this is the reducer we want ...
        if (this.reducers[i].name == name) {
            // ... append this subcriber to the reducer's list of subscribers
            this.reducers[i].subscribers = this.reducers[i].subscribers.concat([method]);
        }
    }
};

Graydux.prototype.dispatch = function (action, data) {
    var i, j;

    // Call all reducers
    var reducerCount = this.reducers.length;
    for (i = 0; i < reducerCount; i++) {
        var reducer = this.reducers[i];
        var keyList = reducer.keyList;

        try {
            // Try to access the parts of the store which the reducer manages.
            // If it fails, send the error to onReducerFail and continue
            this.getState(keyList);
        } catch(error) {
            this.onReducerFail({
                error: error,
                reducer: reducer
            });
            continue;
        }

        // Get new state according to the reducer
        var newState = reducer.method(this.getState(keyList), action, data);
        // Set new state
        this.setState(keyList, newState);

        // Wrap sending new state to subscribers async
        var castNewStateToSubs = function (sub, state) {
            window.setTimeout(function () {
                sub(state);
            });
        }

        // Call all subscribers of each reducer
        var subscriberCount = reducer.subscribers.length;
        for (j = 0; j < subscriberCount; j++) {
            // Send new state async
            castNewStateToSubs(reducer.subscribers[j], newState);
        }
    }
}

Graydux.prototype.getState = function (keyList) {
    if (keyList.length == 0) {
        // getState([]) returns entire store.
        return this.store;
    } else {
        // Otherwise, work our way down the store
        return this._getState(this.store, keyList);
    }

}

Graydux.prototype._getState = function (subStore, keyList) {
    if (Object.keys(subStore).includes(keyList[0])) {
        if (keyList.length == 1) {
            // Base case
            return subStore[keyList[0]];
        } else {
            // Recursive case
            var head = keyList[0];
            var tail = keyList.slice(1);
            return this._getState(subStore[head], tail);
        }
    } else {
        throw {
            name: "getStateError",
            message: "Invalid key: " + keyList[0]
        };
    }
}

Graydux.prototype.setState = function (keyList, data) {
    if (keyList.length == 0) {
        // Passing an empty list to setState replaces store entirely
        this.store = data;
    } else {
        // Otherwise, just update state at the store subcomponent defined by keyList
        this.store = this._setState(this.store, keyList, data);
    }
}

Graydux.prototype._setState = function (storeSubComponent, keyList, data) {
    if (Object.keys(storeSubComponent).includes(keyList[0])) {
        if (keyList.length == 1) {
            // Base case
            storeSubComponent[keyList[0]] = data;
            return storeSubComponent;
        } else {
            // Recursive case
            var head = keyList[0];
            var tail = keyList.slice(1);
            storeSubComponent[head] = this._setState(storeSubComponent[head], tail, data);
            return storeSubComponent;
        }
    } else {
        throw {
            name: "setStateError",
            message: "Invalid key: " + keyList[0]
        };
    }
}
