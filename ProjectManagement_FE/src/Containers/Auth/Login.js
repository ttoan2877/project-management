import React, { useState, useCallback, useEffect, useContext } from 'react'

import {
  StyleSheet,
  View,
  Text,
  TextInput,
  TouchableOpacity,
} from 'react-native'

import AsyncStorage from '@react-native-async-storage/async-storage'

import { Auth } from '../../Services/AuthService'

import { StoreContext } from '../../App'

const Login = ({ navigation }) => {
  const { setLoading, setCurrentUser } = useContext(StoreContext)

  const [username, setUsername] = useState('lqhuy1')

  const [password, setPassword] = useState('123456')

  const onLogin = useCallback(async () => {
    await setLoading(true)
    const user = await Auth({ username, password })
    await setLoading(false)
    user && (await setCurrentUser(user))
    onNavigate()
  }, [onNavigate, password, setCurrentUser, setLoading, username])

  useEffect(() => {
    const checkToken = async () => {
      const value = await AsyncStorage.getItem('@accessToken')
      if (value) {
        onNavigate()
      }
    }
    checkToken()
  }, [onNavigate])

  const onNavigate = useCallback(() => {
    navigation.navigate('SelectProject')
    navigation.reset({
      index: 0,
      routes: [{ name: 'SelectProject' }],
    })
  }, [navigation])

  return (
    <View style={styles.container}>
      <Text style={styles.brand}>ONE TECH ASIA PROJECT MANAGEMENT</Text>
      <View style={styles.area}>
        <View style={styles.labelInput}>
          <Text style={styles.label}>Username:</Text>
          <TextInput
            value={username}
            style={styles.input}
            onChangeText={(val) => setUsername(val)}
          />
        </View>
        <View style={styles.labelInput}>
          <Text style={styles.label}>Password:</Text>
          <TextInput
            value={password}
            secureTextEntry={true}
            style={styles.input}
            onChangeText={(val) => setPassword(val)}
          />
        </View>
        <TouchableOpacity onPress={onLogin} style={styles.btn}>
          <Text style={styles.btnLabel}>LOGIN</Text>
        </TouchableOpacity>
      </View>
    </View>
  )
}

export default Login

const styles = StyleSheet.create({
  brand: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 32,
    color: 'firebrick',
  },
  header: {
    fontSize: 14,
    fontWeight: 'bold',
    color: 'azure',
  },
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'black',
  },
  area: {
    width: 300,
    height: 250,
    borderRadius: 8,
    padding: 16,
    backgroundColor: 'firebrick',
    alignItems: 'center',
  },
  labelInput: {
    height: 50,
    width: '100%',
    backgroundColor: 'azure',
    borderRadius: 8,
    marginVertical: 16,
    flexDirection: 'row',
    alignItems: 'center',
    padding: 8,
  },
  label: {
    fontSize: 12,
    color: 'firebrick',
    fontWeight: 'bold',
  },
  input: {
    flex: 1,
    marginLeft: 8,
    fontSize: 14,
    borderBottomColor: 'firebrick',
    borderBottomWidth: 1,
  },
  btn: {
    width: 150,
    backgroundColor: 'crimson',
    padding: 8,
    marginVertical: 16,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
  btnLabel: {
    fontSize: 16,
    color: 'azure',
    fontWeight: 'bold',
  },
})
