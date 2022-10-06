import React from 'react';
import Navbar from './Navbar';
import PageTitle from '../content/PageTitle';

const style = {
  backgroundImage: `url(${process.env.PUBLIC_URL + '/hero.png'})`,
  backgroundRepeat: 'no-repeat',
  backgroundPosition: 'right',
  backgroundSize: '853px 480px',
  minHeight: '480px',
  width: '100%',
}

const hwrap = {
  backgroundImage: `url(${process.env.PUBLIC_URL + '/wave.png'})`,
  backgroundRepeat: 'x',
  backgroundPosition: 'right',
  minHeight: '480px',
  width: '100%',
}

export default function Header() {
    return (
        <div style={hwrap}>
          <header style={style}>
              <div>
                  <Navbar />
                  <PageTitle />
              </div>
          </header>
        </div>
    )
}