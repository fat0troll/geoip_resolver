package mogrus

// Mogrus is a logger package built on top of Logrus with extended
// (and kinda simplified) approach for multilogging.
//
// Example usage:
//
//     TBF
//
// Output functions are splitted in three things:
//   * ``Debug`` - simple debug printing
//   * ``Debugf`` - debug printing with line formatting
//   * ``Debugln`` - debug printing, same as Debug().
//
// It will try to resemble Logrus as much as possible.
//
// About functions.
//
// Info(f,ln) and Print(f,ln) doing same thing - logging to
// INFO log level.
//
// Be careful while using Fatal(f,ln) functions, because
// it will log only to first logger and then will os.Exit(1)!
// There is no guarantee what logger will be first ATM.
