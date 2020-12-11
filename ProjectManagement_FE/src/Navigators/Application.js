import React, { useContext } from 'react'

import { Text, TouchableOpacity, StyleSheet } from 'react-native'

import { createStackNavigator } from '@react-navigation/stack'
import { NavigationContainer, DrawerActions } from '@react-navigation/native'

import Login from '../Containers/Auth/Login'
import SelectProject from '../Containers/Auth/SelectProject'
import CreateProject from '../Containers/Auth/CreateProject'

import MainNavigator from './Main'
import ModalNavigator from './Modal'

import { StoreContext } from '../App'

const Stack = createStackNavigator()

const ApplicationNavigator = () => {
  const { currentProject } = useContext(StoreContext)
  return (
    <NavigationContainer>
      <Stack.Navigator
        initialRouteName="Login"
        screenOptions={({ navigation }) => ({
          title: currentProject?.name?.toUpperCase() || 'PROJECT',
          headerTitleStyle: styles.headerTitleStyle,
          headerStyle: styles.headerStyle,
          headerLeft: () => (
            <TouchableOpacity
              onPress={() => navigation.dispatch(DrawerActions.toggleDrawer())}
              style={styles.btn}
            >
              <Text style={styles.label}>Menu</Text>
            </TouchableOpacity>
          ),
        })}
      >
        <Stack.Screen
          options={{ headerShown: false }}
          name="Login"
          component={Login}
        />
        <Stack.Screen
          name="SelectProject"
          component={SelectProject}
          options={{ headerShown: false }}
        />
        <Stack.Screen
          name="CreateProject"
          component={CreateProject}
          options={{ headerShown: false }}
        />
        <Stack.Screen name="Main" component={MainNavigator} />
        <Stack.Screen
          options={{ headerShown: false }}
          name="Modal"
          component={ModalNavigator}
        />
      </Stack.Navigator>
    </NavigationContainer>
  )
}

export default ApplicationNavigator

const styles = StyleSheet.create({
  headerTitleStyle: {
    fontWeight: 'bold',
    color: 'white',
    alignSelf: 'center',
  },
  headerStyle: {
    backgroundColor: 'firebrick',
  },
  btn: {
    marginLeft: 16,
  },
  label: {
    color: 'white',
    fontWeight: 'bold',
  },
})
