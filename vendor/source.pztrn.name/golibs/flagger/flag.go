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

// This structure represents addable flag for Flagger.
type Flag struct {
    // Flag name. It will be accessible using this name later.
    Name string
    // Description for help output.
    Description string
    // Type can be one of "bool", "int", "string".
    Type string
    // This value will be reflected.
    DefaultValue interface{}
}
