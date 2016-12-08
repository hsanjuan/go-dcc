// Package dcc implements the DCC protocol for controlling model trains.
// It can support a number of different encoders, which are in charge of
// translating DCC packages into electrical signals. By default, a Raspberry
// Pi driver is provided.
//
// The implementation follows the S-91 Electrical Standard (http://www.nmra.org/sites/default/files/standards/sandrp/pdf/s-9.1_electrical_standards_2006.pdf), the S-92 DCC Communications Standard (http://www.nmra.org/sites/default/files/s-92-2004-07.pdf) and the S-9.2.1 Extended Packet Formats for Digital Command Control standard (http://www.nmra.org/sites/default/files/s-9.2.1_2012_07.pdf).
package dcc
