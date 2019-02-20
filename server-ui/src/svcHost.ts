export const grpcHost = window.location.port === '3000' ?
    `http://172.23.225.30:30090` :
    `${window.location.protocol}//${window.location.host}`;