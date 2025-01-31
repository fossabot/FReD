package peering

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is an peering client to communicate with peers.
type Client struct {
	conn        map[string]peering.NodeClient
	credentials credentials.TransportCredentials
}

// NewClient creates a new empty client to communicate with peers.
func NewClient(certFile string, keyFile string, caFile string) *Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")

		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs, err := x509.SystemCertPool()

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load root certificates")
		return nil
	}

	loaded, err := ioutil.ReadFile(caFile)

	if err != nil {
		log.Fatal().Msgf("unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	return &Client{
		conn:        make(map[string]peering.NodeClient),
		credentials: credentials.NewTLS(tlsConfig),
	}
}

// getClient creates a new connection to a server or uses an existing one.
func (c *Client) getClient(host string) (peering.NodeClient, error) {
	if client, ok := c.conn[host]; ok {
		return client, nil
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(c.credentials))

	if err != nil {
		log.Error().Err(err).Msg("Cannot create Grpc connection")
		return nil, errors.New(err)
	}

	log.Debug().Msgf("Interclient: Created Connection to %s", host)

	client := peering.NewNodeClient(conn)
	c.conn[host] = client
	return client, nil
}

// Destroy currently does nothing, but might delete open connections if they are implemented
func (c *Client) Destroy() {
}

// SendCreateKeygroup sends this command to the server at this address
func (c *Client) SendCreateKeygroup(host string, kgname fred.KeygroupName, expiry int) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.CreateKeygroup(context.Background(), &peering.CreateKeygroupRequest{Keygroup: string(kgname), Expiry: int64(expiry)})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendDeleteKeygroup sends this command to the server at this address
func (c *Client) SendDeleteKeygroup(host string, kgname fred.KeygroupName) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.DeleteKeygroup(context.Background(), &peering.DeleteKeygroupRequest{Keygroup: string(kgname)})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendUpdate sends this command to the server at this address
func (c *Client) SendUpdate(host string, kgname fred.KeygroupName, id string, value string) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.PutItem(context.Background(), &peering.PutItemRequest{
		Keygroup: string(kgname),
		Id:       id,
		Data:     value,
	})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendAppend sends this command to the server at this address
func (c *Client) SendAppend(host string, kgname fred.KeygroupName, id string, value string) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.AppendItem(context.Background(), &peering.AppendItemRequest{
		Keygroup: string(kgname),
		Id:       id,
		Data:     value,
	})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendDelete sends this command to the server at this address
func (c *Client) SendDelete(host string, kgname fred.KeygroupName, id string) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.DeleteItem(context.Background(), &peering.DeleteItemRequest{
		Keygroup: string(kgname),
		Id:       id,
	})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendAddReplica sends this command to the server at this address
func (c *Client) SendAddReplica(host string, kgname fred.KeygroupName, node fred.Node, expiry int) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.AddReplica(context.Background(), &peering.AddReplicaRequest{
		NodeId:   string(node.ID),
		Keygroup: string(kgname),
		Expiry:   int64(expiry),
	})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendRemoveReplica sends this command to the server at this address
func (c *Client) SendRemoveReplica(host string, kgname fred.KeygroupName, node fred.Node) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.RemoveReplica(context.Background(), &peering.RemoveReplicaRequest{
		NodeId:   string(node.ID),
		Keygroup: string(kgname),
	})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendGetItem sends this command to the server at this address
func (c *Client) SendGetItem(host string, kgname fred.KeygroupName, id string) (fred.Item, error) {
	client, err := c.getClient(host)

	if err != nil {
		return fred.Item{}, errors.New(err)
	}

	res, err := client.GetItem(context.Background(), &peering.GetItemRequest{
		Keygroup: string(kgname),
		Id:       id,
	})

	if err != nil {
		return fred.Item{}, errors.New(err)
	}

	return fred.Item{
		Keygroup: kgname,
		ID:       id,
		Val:      res.Data,
	}, nil
}

// SendGetAllItems sends this command to the server at this address
func (c *Client) SendGetAllItems(host string, kgname fred.KeygroupName) ([]fred.Item, error) {
	client, err := c.getClient(host)

	if err != nil {
		return nil, err
	}

	res, err := client.GetAllItems(context.Background(), &peering.GetAllItemsRequest{
		Keygroup: string(kgname),
	})

	if err != nil {
		return nil, errors.New(err)
	}

	d := make([]fred.Item, len(res.Data))

	for i, item := range res.Data {
		d[i] = fred.Item{
			Keygroup: kgname,
			ID:       item.Id,
			Val:      item.Data,
		}
	}

	return d, nil
}
