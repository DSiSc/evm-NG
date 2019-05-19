pragma solidity ^0.4.0;

// SolidityBuffer used to cache solidity temporary data
contract SolidityBuffer {
    bytes private buffer;
    // read data from the buffer
    function read(uint pos, uint size) view public returns (bytes) {
        bytes memory subdata;
        if (pos >= buffer.length) {
            return subdata;
        } else if ((pos + size) < buffer.length) {
            subdata = new bytes(size);
        } else {
            subdata = new bytes(buffer.length - pos);
        }
        for (uint160 i = 0; i < subdata.length; i++) {
            subdata[i] = buffer[pos + i];
        }
        return subdata;
    }

    // write data to buffer
    function write(bytes content) public {
        for (uint i = 0; i < content.length; i++) {
            buffer.push(content[i]);
        }
    }

    // get total length of the data in buffer
    function length() view public returns (uint) {
        return buffer.length;
    }

    // close the buffer
    function close() public {
        delete buffer;
    }
}