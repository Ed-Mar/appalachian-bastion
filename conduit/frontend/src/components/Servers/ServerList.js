import React from 'react';
import axios from 'axios';


class ServerList extends React.Component {
    state = {
        servers: []
    }

    async componentDidMount() {
        try {
            const response = await axios.get('http://localhost:9090/servers');
            console.log(response);
            let servers = response.data
            this.setState({servers})
        } catch (error) {
            console.error(error);
        }
    }

    //
    // componentDidMount() {
    //     axios.get(`http://localhost:9090/servers`)
    //         .then(res => {
    //             const servers = res.data;
    //             this.setState({ servers });
    //         })
    // }


    render() {
        return (
            <div>
                {this.state.servers.map(({ id , name, description  }) => (
                    <p key={id}>ID:{id} | Name: {name} | Description: {description}</p>
                ))}
            </div>


        )
    }
}
export default ServerList

