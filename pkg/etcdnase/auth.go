package etcdnase

import (
	"fmt"
	"strings"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// RevokeUserPermissions removes user's permission to perform method on kg by deleting the key in etcd.
func (n *NameService) RevokeUserPermissions(user string, method fred.Method, kg fred.KeygroupName) error {
	prefix := fmt.Sprintf(fmtUserPermissionStringPrefix, user, string(kg))
	return n.delete(prefix+string(method), prefix)
}

// AddUserPermissions adds user's permission to perform method on kg by adding the key to etcd.
func (n *NameService) AddUserPermissions(user string, method fred.Method, kg fred.KeygroupName) error {
	prefix := fmt.Sprintf(fmtUserPermissionStringPrefix, user, string(kg))
	return n.put(prefix+string(method), "ok", prefix)
}

// GetUserPermissions returns a set of all of the user's permissions on kg from etcd.
func (n *NameService) GetUserPermissions(user string, kg fred.KeygroupName) (map[fred.Method]struct{}, error) {
	res, err := n.getPrefix(fmt.Sprintf(fmtUserPermissionStringPrefix, user, string(kg)))

	if err != nil {
		return nil, err
	}

	permissions := make(map[fred.Method]struct{})

	for k := range res {
		permissions[fred.Method(strings.Split(k, sep)[5])] = struct{}{}
	}

	return permissions, nil
}
