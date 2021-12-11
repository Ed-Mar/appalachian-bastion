import Head from 'next/head'

import styles from '../styles/Home.module.css'
import Link from 'next/link'
import Layout from '../components/layout'



export default function Home() {
  return (
    <Layout>
      <Head>
        <title>Conduit Landing</title>
      </Head>
      <p>Landing Page</p>
      <Link href="/servers">
        <a>
        <p>Servers</p>
        </a>
      </Link>
      
    </Layout>
  )
}