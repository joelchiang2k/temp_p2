import Head from 'next/head'
import Image from 'next/image'
import { Inter } from 'next/font/google'
import styles from '@/styles/Home.module.css'
import { useState, useEffect } from 'react'
import FileInput from './zipUpload'
import { output } from '../../next.config.cjs'

const inter = Inter({ subsets: ['latin'] })

/*export async function getStaticProps(context){
  const res = await fetch('http://localhost:8000/package')
const message = await res.json();

  return { 
    props: {message}
    }  ;
}*/

export default function Home({message}) {

  function handleClick() {
    alert('clicked!');
  }

  const [url, setURL] = useState('');
  /*const [packageName, setName] = useState('');
  const [packageVersion, setVersion] = useState('');*/
  const [postID, setPostID] = useState();
  const [zipData, setZipData] = useState('');
  const [postName, setPostName] = useState('');
  const [postVersion, setPostVersion] = useState('');
  const [rows, setRows] = useState([]);
  const [RegEx, setRegex] = useState('');
  

  useEffect(() => {
    const storeRows = JSON.parse(localStorage.getItem('rows'));
    if(storeRows){
      setRows(storeRows);
    }
  }, []);

  useEffect(() => {
    localStorage.setItem('rows', JSON.stringify(rows));
  }, [rows]);

  const addRow = () => {
    const newRow = { ID: postID, Name: postName, Version: postVersion };
    setRows([...rows, newRow]);
    setPostID('');
    setPostName('');
    setPostVersion('');
  };

  const handleZipUpload = (responseJSON) => {
    const addRow = {ID: responseJSON.metadata.ID, Name: responseJSON.metadata.Name, Version: responseJSON.metadata.Version};
    setRows([...rows, addRow]);
  }

  const handleFileSubmit = (base64Data) => {
    //setFileData(base64Data)
    setZipData(base64Data)
    //e.preventDefault();
    //sendPostRequest(b64String);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    //const dataStruct = { packageName, packageVersion, url}; //content};
    const dataStruct = { url };
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}` + "/package", {
      method: 'POST',
      headers: { "Content-Type": "application/json",
                  "X-Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

      },
      body: JSON.stringify(dataStruct)
    }).then(response => response.json())
      .then(responseJSON => {
        console.log(responseJSON);
        console.log(responseJSON.metadata)
        const addRow = {ID: responseJSON.metadata.ID, Name: responseJSON.metadata.Name, Version: responseJSON.metadata.Version};
        setRows([...rows, addRow]);

      }).catch(error => console.error(error));
    
    //const addRow = {ID: response.metadata.id, Name: response.metadata.id, Version: response.metadata.Version};
    //setRows([...rows, addRow]);
    //const addRow = { ID: response.}
  }

  const handleReset = () => {
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}` + "/reset", {
      method: 'DELETE',
      headers: { "Content-Type": "application/json",
                  "X-Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

      }
    }).then(response => {
      if(!response.ok){
        throw new Error("Failed to reset database");
      }
      setRows([]);
    }).catch(error => console.error(error));
  }

  const handleDeleteRow = (rowID) => {
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}` + "/package/" + String(rowID), {
      method: 'DELETE',
      headers: { "Content-Type": "application/json",
                  "X-Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
      }
    })
    .then(response => response.json())
    .then(() => {
      const newRows = rows.filter(row => row.ID !== rowID);
      setRows(newRows)
    })
    .catch(error => console.error(error));
  }

  const handleRateRow = (rowID) => {
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}` + "/package/" + String(rowID) + "/rate", {
      method: 'GET',
      headers: { "Content-Type": "application/json",
                  "X-Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
      }
    })
    .then(response => response.json())
    //.then(data => window.alert(JSON.stringify(data)))
    .then(data => {
      const windowObj = document.createElement('div');

      windowObj.style.display = 'flex';
      windowObj.style.flexDirection = 'column';
      windowObj.style.alignItems = 'center';
      windowObj.style.justifyContent = 'center';
      windowObj.style.position = 'fixed';
      windowObj.style.zIndex = '9999';
      windowObj.style.top = '0';
      windowObj.style.left = '0';
      windowObj.style.width = '100vw';
      windowObj.style.height = '100vh';
      windowObj.style.backgroundColor = 'rgba(0, 0, 0, 0.5)';

      const windowContentObj = document.createElement('div');
      windowContentObj.style.display = 'flex';
      windowContentObj.style.flexDirection = 'column';
      windowContentObj.style.alignItems = 'center';
      windowContentObj.style.justifyContent = 'center';
      windowContentObj.style.width = '50%';
      windowContentObj.style.backgroundColor = 'blue';
      windowContentObj.style.padding = '2rem';

      const headerObj = document.createElement('h2');
      headerObj.innerText = "Metrics for Package with ID:" + String(rowID);

      const outputObj = document.createElement('textarea');
      outputObj.value = JSON.stringify(data, null, 2);
      outputObj.style.width = '100%';
      outputObj.style.minHeight = '10rem';
      outputObj.style.resize = 'none';
      outputObj.style.padding = '0.5rem';
      
      windowContentObj.appendChild(headerObj);
      windowContentObj.appendChild(outputObj);
      
      const xButton = document.createElement('button');
      xButton.innerText = "x";
      xButton.addEventListener('click', () => {
        windowObj.remove();
      })
      
      windowContentObj.appendChild(xButton);

      windowObj.appendChild(windowContentObj);
      document.body.appendChild(windowObj);
    })
    .catch(error => console.error(error));
  }

  const handleRegex = (e) => {
    const dataStruct = { RegEx };
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}` + "/package" + "/byRegEx", {
      method: 'POST',
      headers: { "Content-Type": "application/json",
                  "X-Authorization": "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
      },
      body: JSON.stringify(dataStruct)
    })
    .then(response => response.json())
    .then(data => window.alert(JSON.stringify(data)))
  }
  /*const sendPostRequest = (b64String) => {
      fetch(`${process.env.NEXT_PUBLIC_BACKEND_API}`, {
          method: 'POST',
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ Content: b64String})
      })
  };//response?*/

  /*function FileInput() {
    const [selectedZip, setSelectedZip] = useState(null)

  }*/
  return (
    <>
      <Head>
        <title>Create Next App</title>
        <meta name="description" content="Generated by create next app" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <div className={styles.center}>
          <h1> Welcome to 461 Part 2! </h1>
        </div>
  {/*      <div>
          <a href="http://localhost:8000/package" target="_blank">
            <button> Sample API button </button>
          </a>
  </div>*/}
        {/*<div>message: {message.message}</div>*/}
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
            <div class="flex-child" style={{ flexGrow: 1}}>
              <h2 style={{paddingBottom: '15px', paddingLeft: '200px'}}>Package Ingestion From URL</h2>
              <form onSubmit={handleSubmit}>
                <label style={{paddingRight: '15px', paddingLeft: '200px'}}>Enter URL</label>
                <input
                  type="text" 
                  required
                  value={url}
                  onChange={(e) => setURL(e.target.value)}
                  />
                <button>Submit</button>
              </form>
            </div>
            <div class="flex-child" style={{ flexGrow: 1, paddingLeft: '100px'}}>
              <h2 style={{paddingBottom: '15px'}}>Upload Package</h2>
              <FileInput handleSubmit={handleFileSubmit} onZipUpload={handleZipUpload}/>
            </div>
        </div>

        <div style={{paddingTop: '100px', paddingLeft: '200px'}}>
          <h2>Package Search</h2> 
          <form onSubmit={handleRegex}>
            <label>Package Name:</label>
            <input
              type="text"
              required
              value={RegEx}
              onChange={(e) => setRegex(e.target.value)}
              />
            <button>Search</button>
          </form>
        </div>
        <div style={{display: 'flex', justifyContent:'flex-end', paddingRight: '100px', paddingTop: '100px'}}>
          <button onClick={handleReset}>Reset Database</button>
        </div>

        <center style={{paddingTop: '75px'}}>
          <h1 style={{paddingBottom: '15px'}}>Packages</h1>
          <table>
            <thead>
              <tr>
                <th style={{paddingRight: '150px'}}>ID</th>
                <th style={{paddingRight: '150px'}}>Name</th> 
                <th>Version</th>
              </tr>
            </thead>
              {rows.map((row) => (
                <tr key={row.ID}>
                  <td>{row.ID}</td>
                  <td>{row.Name}</td>
                  <td>{row.Version}</td>
                  <td>
                    <button onClick={() => handleRateRow(row.ID)}>Rate</button>
                  </td>
                  <td>
                    <button>Update</button>
                  </td>
                  <td>
                    <button onClick={() => handleDeleteRow(row.ID) }>Delete</button>
                  </td>
                </tr>
              ))}
            <tbody>

            </tbody>
          </table>
        </center>
      </main>
    </>
  )
}