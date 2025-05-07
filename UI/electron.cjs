const { app, BrowserWindow } = require('electron')
const path = require('path')

function createWindow () {
  const win = new BrowserWindow({
    width: 1920,
    height: 1080,
    webPreferences: {
      contextIsolation: true,
      nodeIntegration: true,
      webSecurity: false
    }
  })

  // Carga la app Vue en modo desarrollo o producción
  const indexPath = path.join(__dirname, 'dist', 'index.html')
  console.log('Loading file from:', indexPath)
  
  // Verifica si el archivo existe
  const fs = require('fs')
  if (fs.existsSync(indexPath)) {
    win.loadFile(indexPath)
  } else {
    console.error('index.html not found at:', indexPath)
  }

  // Abre las DevTools en desarrollo
  win.webContents.openDevTools()
}

app.whenReady().then(createWindow)

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow()
  }
}) 