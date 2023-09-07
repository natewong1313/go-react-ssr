"use strict";
(() => {
  var z = Object.create;
  var O = Object.defineProperty;
  var W = Object.getOwnPropertyDescriptor;
  var Y = Object.getOwnPropertyNames;
  var G = Object.getPrototypeOf,
    J = Object.prototype.hasOwnProperty;
  var _ = (e, t) => () => (t || e((t = { exports: {} }).exports, t), t.exports);
  var K = (e, t, r, n) => {
    if ((t && typeof t == "object") || typeof t == "function")
      for (let o of Y(t))
        !J.call(e, o) &&
          o !== r &&
          O(e, o, {
            get: () => t[o],
            enumerable: !(n = W(t, o)) || n.enumerable,
          });
    return e;
  };
  var E = (e, t, r) => (
    (r = e != null ? z(G(e)) : {}),
    K(
      t || !e || !e.__esModule
        ? O(r, "default", { value: e, enumerable: !0 })
        : r,
      e
    )
  );
  var L = _((u) => {
    "use strict";
    var d = Symbol.for("react.element"),
      Q = Symbol.for("react.portal"),
      X = Symbol.for("react.fragment"),
      Z = Symbol.for("react.strict_mode"),
      ee = Symbol.for("react.profiler"),
      te = Symbol.for("react.provider"),
      re = Symbol.for("react.context"),
      ne = Symbol.for("react.forward_ref"),
      oe = Symbol.for("react.suspense"),
      ue = Symbol.for("react.memo"),
      se = Symbol.for("react.lazy"),
      x = Symbol.iterator;
    function ce(e) {
      return e === null || typeof e != "object"
        ? null
        : ((e = (x && e[x]) || e["@@iterator"]),
          typeof e == "function" ? e : null);
    }
    var N = {
        isMounted: function () {
          return !1;
        },
        enqueueForceUpdate: function () {},
        enqueueReplaceState: function () {},
        enqueueSetState: function () {},
      },
      I = Object.assign,
      q = {};
    function y(e, t, r) {
      (this.props = e),
        (this.context = t),
        (this.refs = q),
        (this.updater = r || N);
    }
    y.prototype.isReactComponent = {};
    y.prototype.setState = function (e, t) {
      if (typeof e != "object" && typeof e != "function" && e != null)
        throw Error(
          "setState(...): takes an object of state variables to update or a function which returns an object of state variables."
        );
      this.updater.enqueueSetState(this, e, t, "setState");
    };
    y.prototype.forceUpdate = function (e) {
      this.updater.enqueueForceUpdate(this, e, "forceUpdate");
    };
    function T() {}
    T.prototype = y.prototype;
    function R(e, t, r) {
      (this.props = e),
        (this.context = t),
        (this.refs = q),
        (this.updater = r || N);
    }
    var b = (R.prototype = new T());
    b.constructor = R;
    I(b, y.prototype);
    b.isPureReactComponent = !0;
    var g = Array.isArray,
      A = Object.prototype.hasOwnProperty,
      w = { current: null },
      D = { key: !0, ref: !0, __self: !0, __source: !0 };
    function V(e, t, r) {
      var n,
        o = {},
        s = null,
        c = null;
      if (t != null)
        for (n in (t.ref !== void 0 && (c = t.ref),
        t.key !== void 0 && (s = "" + t.key),
        t))
          A.call(t, n) && !D.hasOwnProperty(n) && (o[n] = t[n]);
      var f = arguments.length - 2;
      if (f === 1) o.children = r;
      else if (1 < f) {
        for (var i = Array(f), a = 0; a < f; a++) i[a] = arguments[a + 2];
        o.children = i;
      }
      if (e && e.defaultProps)
        for (n in ((f = e.defaultProps), f)) o[n] === void 0 && (o[n] = f[n]);
      return {
        $$typeof: d,
        type: e,
        key: s,
        ref: c,
        props: o,
        _owner: w.current,
      };
    }
    function ie(e, t) {
      return {
        $$typeof: d,
        type: e.type,
        key: t,
        ref: e.ref,
        props: e.props,
        _owner: e._owner,
      };
    }
    function C(e) {
      return typeof e == "object" && e !== null && e.$$typeof === d;
    }
    function fe(e) {
      var t = { "=": "=0", ":": "=2" };
      return (
        "$" +
        e.replace(/[=:]/g, function (r) {
          return t[r];
        })
      );
    }
    var P = /\/+/g;
    function k(e, t) {
      return typeof e == "object" && e !== null && e.key != null
        ? fe("" + e.key)
        : t.toString(36);
    }
    function h(e, t, r, n, o) {
      var s = typeof e;
      (s === "undefined" || s === "boolean") && (e = null);
      var c = !1;
      if (e === null) c = !0;
      else
        switch (s) {
          case "string":
          case "number":
            c = !0;
            break;
          case "object":
            switch (e.$$typeof) {
              case d:
              case Q:
                c = !0;
            }
        }
      if (c)
        return (
          (c = e),
          (o = o(c)),
          (e = n === "" ? "." + k(c, 0) : n),
          g(o)
            ? ((r = ""),
              e != null && (r = e.replace(P, "$&/") + "/"),
              h(o, t, r, "", function (a) {
                return a;
              }))
            : o != null &&
              (C(o) &&
                (o = ie(
                  o,
                  r +
                    (!o.key || (c && c.key === o.key)
                      ? ""
                      : ("" + o.key).replace(P, "$&/") + "/") +
                    e
                )),
              t.push(o)),
          1
        );
      if (((c = 0), (n = n === "" ? "." : n + ":"), g(e)))
        for (var f = 0; f < e.length; f++) {
          s = e[f];
          var i = n + k(s, f);
          c += h(s, t, r, i, o);
        }
      else if (((i = ce(e)), typeof i == "function"))
        for (e = i.call(e), f = 0; !(s = e.next()).done; )
          (s = s.value), (i = n + k(s, f++)), (c += h(s, t, r, i, o));
      else if (s === "object")
        throw (
          ((t = String(e)),
          Error(
            "Objects are not valid as a React child (found: " +
              (t === "[object Object]"
                ? "object with keys {" + Object.keys(e).join(", ") + "}"
                : t) +
              "). If you meant to render a collection of children, use an array instead."
          ))
        );
      return c;
    }
    function v(e, t, r) {
      if (e == null) return e;
      var n = [],
        o = 0;
      return (
        h(e, n, "", "", function (s) {
          return t.call(r, s, o++);
        }),
        n
      );
    }
    function le(e) {
      if (e._status === -1) {
        var t = e._result;
        (t = t()),
          t.then(
            function (r) {
              (e._status === 0 || e._status === -1) &&
                ((e._status = 1), (e._result = r));
            },
            function (r) {
              (e._status === 0 || e._status === -1) &&
                ((e._status = 2), (e._result = r));
            }
          ),
          e._status === -1 && ((e._status = 0), (e._result = t));
      }
      if (e._status === 1) return e._result.default;
      throw e._result;
    }
    var l = { current: null },
      m = { transition: null },
      ae = {
        ReactCurrentDispatcher: l,
        ReactCurrentBatchConfig: m,
        ReactCurrentOwner: w,
      };
    u.Children = {
      map: v,
      forEach: function (e, t, r) {
        v(
          e,
          function () {
            t.apply(this, arguments);
          },
          r
        );
      },
      count: function (e) {
        var t = 0;
        return (
          v(e, function () {
            t++;
          }),
          t
        );
      },
      toArray: function (e) {
        return (
          v(e, function (t) {
            return t;
          }) || []
        );
      },
      only: function (e) {
        if (!C(e))
          throw Error(
            "React.Children.only expected to receive a single React element child."
          );
        return e;
      },
    };
    u.Component = y;
    u.Fragment = X;
    u.Profiler = ee;
    u.PureComponent = R;
    u.StrictMode = Z;
    u.Suspense = oe;
    u.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED = ae;
    u.cloneElement = function (e, t, r) {
      if (e == null)
        throw Error(
          "React.cloneElement(...): The argument must be a React element, but you passed " +
            e +
            "."
        );
      var n = I({}, e.props),
        o = e.key,
        s = e.ref,
        c = e._owner;
      if (t != null) {
        if (
          (t.ref !== void 0 && ((s = t.ref), (c = w.current)),
          t.key !== void 0 && (o = "" + t.key),
          e.type && e.type.defaultProps)
        )
          var f = e.type.defaultProps;
        for (i in t)
          A.call(t, i) &&
            !D.hasOwnProperty(i) &&
            (n[i] = t[i] === void 0 && f !== void 0 ? f[i] : t[i]);
      }
      var i = arguments.length - 2;
      if (i === 1) n.children = r;
      else if (1 < i) {
        f = Array(i);
        for (var a = 0; a < i; a++) f[a] = arguments[a + 2];
        n.children = f;
      }
      return { $$typeof: d, type: e.type, key: o, ref: s, props: n, _owner: c };
    };
    u.createContext = function (e) {
      return (
        (e = {
          $$typeof: re,
          _currentValue: e,
          _currentValue2: e,
          _threadCount: 0,
          Provider: null,
          Consumer: null,
          _defaultValue: null,
          _globalName: null,
        }),
        (e.Provider = { $$typeof: te, _context: e }),
        (e.Consumer = e)
      );
    };
    u.createElement = V;
    u.createFactory = function (e) {
      var t = V.bind(null, e);
      return (t.type = e), t;
    };
    u.createRef = function () {
      return { current: null };
    };
    u.forwardRef = function (e) {
      return { $$typeof: ne, render: e };
    };
    u.isValidElement = C;
    u.lazy = function (e) {
      return { $$typeof: se, _payload: { _status: -1, _result: e }, _init: le };
    };
    u.memo = function (e, t) {
      return { $$typeof: ue, type: e, compare: t === void 0 ? null : t };
    };
    u.startTransition = function (e) {
      var t = m.transition;
      m.transition = {};
      try {
        e();
      } finally {
        m.transition = t;
      }
    };
    u.unstable_act = function () {
      throw Error("act(...) is not supported in production builds of React.");
    };
    u.useCallback = function (e, t) {
      return l.current.useCallback(e, t);
    };
    u.useContext = function (e) {
      return l.current.useContext(e);
    };
    u.useDebugValue = function () {};
    u.useDeferredValue = function (e) {
      return l.current.useDeferredValue(e);
    };
    u.useEffect = function (e, t) {
      return l.current.useEffect(e, t);
    };
    u.useId = function () {
      return l.current.useId();
    };
    u.useImperativeHandle = function (e, t, r) {
      return l.current.useImperativeHandle(e, t, r);
    };
    u.useInsertionEffect = function (e, t) {
      return l.current.useInsertionEffect(e, t);
    };
    u.useLayoutEffect = function (e, t) {
      return l.current.useLayoutEffect(e, t);
    };
    u.useMemo = function (e, t) {
      return l.current.useMemo(e, t);
    };
    u.useReducer = function (e, t, r) {
      return l.current.useReducer(e, t, r);
    };
    u.useRef = function (e) {
      return l.current.useRef(e);
    };
    u.useState = function (e) {
      return l.current.useState(e);
    };
    u.useSyncExternalStore = function (e, t, r) {
      return l.current.useSyncExternalStore(e, t, r);
    };
    u.useTransition = function () {
      return l.current.useTransition();
    };
    u.version = "18.2.0";
  });
  var $ = _((ke, U) => {
    "use strict";
    U.exports = L();
  });
  var M = _((S) => {
    "use strict";
    var pe = $(),
      ye = Symbol.for("react.element"),
      de = Symbol.for("react.fragment"),
      _e = Object.prototype.hasOwnProperty,
      ve =
        pe.__SECRET_INTERNALS_DO_NOT_USE_OR_YOU_WILL_BE_FIRED.ReactCurrentOwner,
      he = { key: !0, ref: !0, __self: !0, __source: !0 };
    function F(e, t, r) {
      var n,
        o = {},
        s = null,
        c = null;
      r !== void 0 && (s = "" + r),
        t.key !== void 0 && (s = "" + t.key),
        t.ref !== void 0 && (c = t.ref);
      for (n in t) _e.call(t, n) && !he.hasOwnProperty(n) && (o[n] = t[n]);
      if (e && e.defaultProps)
        for (n in ((t = e.defaultProps), t)) o[n] === void 0 && (o[n] = t[n]);
      return {
        $$typeof: ye,
        type: e,
        key: s,
        ref: c,
        props: o,
        _owner: ve.current,
      };
    }
    S.Fragment = de;
    S.jsx = F;
    S.jsxs = F;
  });
  var j = _((be, B) => {
    "use strict";
    B.exports = M();
  });
  var H = E($()),
    p = E(j());
  function me() {
    let [e, t] = (0, H.useState)(0);
    return (0, p.jsxs)("div", {
      className: "App",
      children: [
        (0, p.jsx)("div", {
          children: (0, p.jsx)("a", {
            href: "https://reactjs.org",
            target: "_blank",
          }),
        }),
        (0, p.jsx)("h1", { children: "Rspack + React + TypeScript" }),
        (0, p.jsxs)("div", {
          className: "card",
          children: [
            (0, p.jsxs)("button", {
              onClick: () => t((r) => r + 1),
              children: ["count is ", e],
            }),
            (0, p.jsxs)("p", {
              children: [
                "Edit ",
                (0, p.jsx)("code", { children: "src/App.tsx" }),
                " and save to test HMR",
              ],
            }),
          ],
        }),
        (0, p.jsx)("p", {
          className: "read-the-docs",
          children: "Click on the Rspack and React logos to learn more",
        }),
      ],
    });
  }
  var we = me;
})();
/*! Bundled license information:

react/cjs/react.production.min.js:
  (**
   * @license React
   * react.production.min.js
   *
   * Copyright (c) Facebook, Inc. and its affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   *)

react/cjs/react-jsx-runtime.production.min.js:
  (**
   * @license React
   * react-jsx-runtime.production.min.js
   *
   * Copyright (c) Facebook, Inc. and its affiliates.
   *
   * This source code is licensed under the MIT license found in the
   * LICENSE file in the root directory of this source tree.
   *)
*/
