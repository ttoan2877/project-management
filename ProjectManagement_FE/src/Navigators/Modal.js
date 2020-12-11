import React from 'react'
import { StyleSheet, TouchableOpacity, Text } from 'react-native'
import { createStackNavigator } from '@react-navigation/stack'

import SearchUser from '../Containers/Members/SearchUser'
import FormTask from '../Containers/Tasks/FormTask'
import CreateRole from '../Containers/Roles/CreateRole'

const Stack = createStackNavigator()

const ModalNavigator = () => {
  return (
    <Stack.Navigator
      screenOptions={({ navigation }) => ({
        headerTitleStyle: styles.headerTitleStyle,
        headerStyle: styles.headerStyle,
        headerLeft: () => (
          <TouchableOpacity
            onPress={() => navigation.goBack()}
            style={styles.btn}
          >
            <Text style={styles.label}>{'Back'}</Text>
          </TouchableOpacity>
        ),
      })}
    >
      <Stack.Screen
        name="SearchUser"
        component={SearchUser}
        options={{ title: 'Search user' }}
      />
      <Stack.Screen
        name="FormTask"
        component={FormTask}
        options={{ title: 'Create task' }}
      />
      <Stack.Screen
        name="CreateRole"
        component={CreateRole}
        options={{ title: 'Create role' }}
      />
    </Stack.Navigator>
  )
}

export default ModalNavigator

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
