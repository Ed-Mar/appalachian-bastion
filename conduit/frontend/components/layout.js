import Head from 'next/head'
import styles from './layout.module.css'
import Link from 'next/link'

export default function Layout({ children, home }) {
  return (
    <div className={styles.container}>
      <Head>

        <link rel="icon" href="/public/favicon.ico" />
        <meta
          name="Conduit"
          content=""
        />
          <title>Conduit</title>
      </Head>
      <header className={styles.header}>
        <Link href="/">
          <h1>conduit</h1>
        </Link>
      </header>
      <main>{children}</main>
      {!home && (
        <div className={styles.backToHome}>
          <Link href="/">
            <a>← Back to home</a>
          </Link>
        </div>
      )}
    </div>
  )
}