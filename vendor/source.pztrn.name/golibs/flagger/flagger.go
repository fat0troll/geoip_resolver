// Flagger - arbitrary CLI flags parser.
//
// Copyright (c) 2017, Stanislav N. aka pztrn.
// All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package flagger

import (
    // stdlib
    "errors"
    "flag"
    "sync"
)

// Flagger implements (kinda) extended CLI parameters parser. As it
// available from CommonContext, these flags will be available to
// whole application.
//
// It uses reflection to determine what kind of variable we should
// parse or get.
type Flagger struct {
    // Flags that was added by user.
    flags map[string]*Flag
    flagsMutex sync.Mutex

    // Flags that will be passed to flag module.
    flagsBool map[string]*bool
    flagsInt map[string]*int
    flagsString map[string]*string
}

// Adds flag to list of flags we will pass to ``flag`` package.
func (f *Flagger) AddFlag(flag *Flag) error {
    _, present := f.flags[flag.Name]
    if present {
        logger.Fatalln("Cannot add flag '" + flag.Name + "' - already added!")
        return errors.New("Cannot add flag '" + flag.Name + "' - already added!")
    }

    f.flags[flag.Name] = flag
    return nil
}

// This function returns boolean value for flag with given name.
// Returns bool value for flag and nil as error on success
// and false bool plus error with text on error.
func (f *Flagger) GetBoolValue(name string) (bool, error) {
    fl, present := f.flagsBool[name]
    if !present {
        return false, errors.New("No such flag: " + name)
    }
    return (*fl), nil
}

// This function returns integer value for flag with given name.
// Returns integer on success and 0 on error.
func (f *Flagger) GetIntValue(name string) (int, error) {
    fl, present := f.flagsInt[name]
    if !present {
        return 0, errors.New("No such flag: " + name)
    }
    return (*fl), nil
}

// This function returns string value for flag with given name.
// Returns string on success or empty string on error.
func (f *Flagger) GetStringValue(name string) (string, error) {
    fl, present := f.flagsString[name]
    if !present {
        return "", errors.New("No such flag: " + name)
    }
    return (*fl), nil
}

// Flagger initialization.
func (f *Flagger) Initialize() {
    logger.Println("Initializing CLI parameters parser...")

    f.flags = make(map[string]*Flag)

    f.flagsBool = make(map[string]*bool)
    f.flagsInt = make(map[string]*int)
    f.flagsString = make(map[string]*string)
}

// This function adds flags from flags map to flag package and parse
// them. They can be obtained later by calling GetTYPEValue(name),
// where TYPE is one of Bool, Int, String.
func (f *Flagger) Parse() {
    for name, fl := range f.flags {
        if fl.Type == "bool" {
            fdef := fl.DefaultValue.(bool)
            f.flagsBool[name] = &fdef
            flag.BoolVar(&fdef, name, fdef, fl.Description)
        } else if fl.Type == "int" {
            fdef := fl.DefaultValue.(int)
            f.flagsInt[name] = &fdef
            flag.IntVar(&fdef, name, fdef, fl.Description)
        } else if fl.Type == "string" {
            fdef := fl.DefaultValue.(string)
            f.flagsString[name] = &fdef
            flag.StringVar(&fdef, name, fdef, fl.Description)
        }
    }

    logger.Println("Parsing CLI parameters...")
    flag.Parse()
}
