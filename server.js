const ws = require('ws');
const readline = require('readline');
const sirv = require('sirv');
const polka = require('polka');
const compress = require('compression')();
const { LENO_PORT=3000 } = process.env;

module.exports = {
    createServer: () => {
        // Setup readline interface
        const interface = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        // Create a WebSocket server
        const wsServer = new ws.Server({ noServer: true });
        wsServer.on('connection', socket => {
            function sender(line) {
                socket.send(line);
            }
            // Make sure the listener is removed so it doesn't send data to a closet socket
            socket.on('close', () => interface.removeListener('line', sender));
            interface.on('line', sender);
        });
        
        const assets = sirv(require('path').resolve(__dirname, 'public'));
        
        const { server } = polka()
            .use(compress, assets)
            .listen(LENO_PORT, err => {
                if (err) throw err;
                console.log(`🌲️ Leño is running on http://localhost:${LENO_PORT}`);
            });

        server.on('upgrade', (request, socket, head) => {
            wsServer.handleUpgrade(request, socket, head, socket => {
                wsServer.emit('connection', socket, request);
            });
        });

        return server;
    }
};