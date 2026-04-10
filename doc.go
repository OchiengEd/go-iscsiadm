// Package iscsiadm provides a lightweight Go wrapper around the system's
// `iscsiadm` command, intended for use in Container Storage Interface drivers.
//
// It exposes high‑level operations such as Login, Logout, Discover and Sessions,
// while translating iSCSI exit codes into idiomatic Go errors.
package iscsiadm
