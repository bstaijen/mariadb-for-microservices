app.factory('LocalStorage', function () {

    /**
     * Prefix for LocalStorage key.
     * @const
     */
    var STORAGE_PREFIX = 'm_';
    var USER_KEY = "user";
    var TOKEN_KEY = "token";

    return {

        /**
         * Save user in LocalStorage. If data is empty, save an empty array.
         * @param data The user
         */
        setUser: function (data) {
            localStorage.setItem(STORAGE_PREFIX + USER_KEY, JSON.stringify(data || []));
        },

        /**
         * Return saved user
         * @returns {Array} user or empty array
         */
        getUser: function () {
            return JSON.parse(localStorage.getItem(STORAGE_PREFIX + USER_KEY));
        },

        /**
         * Remove user from LocalStorage.
         */
        removeUser: function () {
            localStorage.removeItem(STORAGE_PREFIX + USER_KEY);
        },

        /**
         * Is there an user in the LocalStorage.
         * @returns {boolean} true or false
         */
        hasUser: function() {
            return this.getUser() !== null;
        },

        /**
         * Checks whether there is an token available
         * @returns {boolean} true of false
         */
        hasToken: function () {
            return this.getToken() !== null;
        },

        /**
         * Get token from LocalStorage.
         * @returns {string} token.
         */
        getToken: function () {
            return localStorage.getItem(STORAGE_PREFIX + TOKEN_KEY);
        },

        /**
         * Save token in LocalStorage.
         * @param token Token
         */
        setToken: function (token) {
            if (token && token.length > 0) {
                localStorage.setItem(STORAGE_PREFIX + TOKEN_KEY, token);
            }
        },

        /**
         * Remove token from LocalStorage.
         */
        removeToken: function () {
            localStorage.removeItem(STORAGE_PREFIX + TOKEN_KEY);
        }
    }
});